package ssh

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	sshagent "golang.org/x/crypto/ssh/agent"
)

type agent struct {
	sshAgent            sshagent.ExtendedAgent
	pipelinedConnection io.ReadWriter
	waiters             []chan agentResponse
	lock                sync.Mutex
	signers             []ssh.Signer
	signersByPath       map[string]ssh.Signer
}

func newAgent(pipelineTimeout time.Duration) (*agent, error) {
	sock, err := agentConnection()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to SSH agent: %s", err)
	}
	a := &agent{sshAgent: sshagent.NewClient(sock), waiters: make([]chan agentResponse, 0, 64)}

	if _, ok := sock.(*net.UnixConn); ok {
		// Determine whether we can use the faster pipelined ssh agent protocol
		a.pipelinedConnection, _ = agentConnection()
		go a.readLoop()
		if !a.canDoPipelinedSigning(pipelineTimeout) {
			a.pipelinedConnection.(*net.UnixConn).Close()
			a.pipelinedConnection = nil
			logrus.Warnf("Using slow ssh agent, see https://herd.seveas.net/documentation/ssh_agent.html to fix this")
		}
	}
	a.signers, err = a.Signers()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve keys from SSH agent: %s", err)
	}
	if len(a.signers) == 0 {
		return nil, fmt.Errorf("No keys found in ssh agent")
	}

	a.signersByPath = make(map[string]ssh.Signer)
	for _, signer := range a.signers {
		comment := signer.PublicKey().(*sshagent.Key).Comment
		a.signersByPath[comment] = signer
	}

	return a, nil
}

func agentConnection() (io.ReadWriter, error) {
	if sockPath, ok := os.LookupEnv("SSH_AUTH_SOCK"); ok {
		return net.Dial("unix", sockPath)
	} else if sock := findPageant(); sock != nil {
		return sock, nil
	}
	if _, ok := os.LookupEnv("SSH_CONNECTION"); ok {
		return nil, fmt.Errorf("No ssh agent found in environment, make sure your ssh agent is running and forwarded")
	}
	return nil, fmt.Errorf("No ssh agent found in environment, make sure your ssh agent is running")
}

type agentResponse struct {
	data []byte
	err  error
}

// Test the ssh agent by sending a bunch of requests in a pipelined way. If
// they are not answered within the specified interval (50ms by default), the
// ssh agent is too old and suffers from the bug solved in
// https://github.com/openssh/openssh-portable/pull/183
func (a *agent) canDoPipelinedSigning(timeout time.Duration) bool {
	keys, err := a.List()
	if err != nil || len(keys) == 0 {
		// This is a lie, but avoids double errors: the next step checks
		// whether there even are keys and will throw a better error
		return true
	}
	tests := 10
	c := make(chan bool)
	t := time.NewTicker(timeout)
	defer t.Stop()
	for i := 0; i < tests; i++ {
		go func() { _, err = a.Sign(keys[0], []byte("Test")); c <- (err == nil) }()
	}
	for i := 0; i < tests; i++ {
		select {
		case v := <-c:
			if !v {
				return false
			}
		case <-t.C:
			return false
		}
	}
	return true
}

func (a *agent) readLoop() {
	for {
		data, err := a.readSingleReply()
		a.lock.Lock()
		if len(a.waiters) == 0 {
			// We've hit EOF and readSingleReply returns instant errors
			a.lock.Unlock()
			break
		}
		ch := a.waiters[0]
		a.waiters = a.waiters[1:]
		a.lock.Unlock()
		ch <- agentResponse{data: data, err: err}
		if err != nil {
			break
		}
	}
}

func (a *agent) readSingleReply() ([]byte, error) {
	var respSizeBuf [4]byte
	if _, err := io.ReadFull(a.pipelinedConnection, respSizeBuf[:]); err != nil {
		return nil, err
	}
	respSize := binary.BigEndian.Uint32(respSizeBuf[:])
	buf := make([]byte, respSize)
	if _, err := io.ReadFull(a.pipelinedConnection, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (a *agent) List() ([]*sshagent.Key, error) {
	return a.sshAgent.List()
}

type agentSignRequest struct {
	Key   []byte `sshtype:"13"`
	Data  []byte
	Flags uint32
}

func (a *agent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	if a.pipelinedConnection == nil {
		return a.sshAgent.Sign(key, data)
	}
	req := ssh.Marshal(agentSignRequest{Key: key.Marshal(), Data: data, Flags: uint32(0)})
	if len(req) > math.MaxUint32 {
		return nil, errors.New("ssh agent request too large")
	}
	msg := make([]byte, 4+len(req))
	binary.BigEndian.PutUint32(msg, uint32(len(req))) // nolint:gosec // Overflow is checked above
	copy(msg[4:], req)

	ch := make(chan agentResponse)
	a.lock.Lock()
	_, err := a.pipelinedConnection.Write(msg)
	if err != nil {
		a.lock.Unlock()
		return nil, err
	}
	a.waiters = append(a.waiters, ch)
	a.lock.Unlock()

	resp := <-ch

	if resp.err != nil {
		return nil, resp.err
	}

	if resp.data[0] != 14 {
		return nil, errors.New("ssh agent failed to sign the message")
	}

	var sig ssh.Signature
	if err := ssh.Unmarshal(resp.data[5:], &sig); err != nil {
		return nil, err
	}
	return &sig, nil
}

func (a *agent) Add(key sshagent.AddedKey) error {
	return a.sshAgent.Add(key)
}

func (a *agent) Remove(key ssh.PublicKey) error {
	return a.sshAgent.Remove(key)
}

func (a *agent) RemoveAll() error {
	return a.sshAgent.RemoveAll()
}

func (a *agent) Lock(passphrase []byte) error {
	return a.sshAgent.Lock(passphrase)
}

func (a *agent) Unlock(passphrase []byte) error {
	return a.sshAgent.Unlock(passphrase)
}

func (a *agent) Signers() ([]ssh.Signer, error) {
	signers, err := a.sshAgent.Signers()
	if err != nil {
		return nil, err
	}
	if a.pipelinedConnection == nil {
		return signers, nil
	}

	ret := make([]ssh.Signer, len(signers))
	for i, s := range signers {
		ret[i] = &signer{a, s.PublicKey()}
	}

	return ret, nil
}

func (a *agent) SignWithFlags(key ssh.PublicKey, data []byte, flags sshagent.SignatureFlags) (*ssh.Signature, error) {
	return a.sshAgent.SignWithFlags(key, data, flags)
}

func (a *agent) Extension(extensionType string, contents []byte) ([]byte, error) {
	return a.sshAgent.Extension(extensionType, contents)
}

func (a *agent) SignersForPathCallback(path string) func() ([]ssh.Signer, error) {
	return func() ([]ssh.Signer, error) {
		signers := a.SignersForPath(path)
		if len(signers) == 0 {
			return nil, fmt.Errorf("SSH key %s was not found in the SSH agent", path)
		}
		return signers, nil
	}
}

func (a *agent) SignersForPath(path string) []ssh.Signer {
	if path != "" {
		if k, ok := a.signersByPath[path]; ok {
			return []ssh.Signer{k}
		} else {
			return []ssh.Signer{}
		}
	}
	return a.signers
}

type signer struct {
	agent *agent
	key   ssh.PublicKey
}

func (s *signer) PublicKey() ssh.PublicKey {
	return s.key
}

func (s *signer) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return s.agent.Sign(s.key, data)
}
