package herd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

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

type sshAgentClient struct {
	client  agent.ExtendedAgent
	conn    io.ReadWriter
	waiters []chan agentResponse
	lock    sync.Mutex
}

func NewSshAgentClient(conn, conn2 io.ReadWriter) *sshAgentClient {
	client := &sshAgentClient{client: agent.NewClient(conn), conn: conn2, waiters: make([]chan agentResponse, 0, 64)}
	go client.readLoop()
	return client
}

// Test the ssh agent by sending a bunch of requests in a pipelined way. If
// they are not answered within the specified interval (50ms by default), the
// ssh agent is too old and suffers from the bug solved in
// https://github.com/openssh/openssh-portable/pull/183
func (a *sshAgentClient) functional(timeout time.Duration) bool {
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

func (a *sshAgentClient) readLoop() {
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
	}
}

func (a *sshAgentClient) readSingleReply() (reply []byte, err error) {
	var respSizeBuf [4]byte
	if _, err = io.ReadFull(a.conn, respSizeBuf[:]); err != nil {
		return nil, err
	}
	respSize := binary.BigEndian.Uint32(respSizeBuf[:])
	buf := make([]byte, respSize)
	if _, err = io.ReadFull(a.conn, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (a *sshAgentClient) List() ([]*agent.Key, error) {
	return a.client.List()
}

type agentSignRequest struct {
	Key   []byte `sshtype:"13"`
	Data  []byte
	Flags uint32
}

func (a *sshAgentClient) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	req := ssh.Marshal(agentSignRequest{Key: key.Marshal(), Data: data, Flags: uint32(0)})
	msg := make([]byte, 4+len(req))
	binary.BigEndian.PutUint32(msg, uint32(len(req)))
	copy(msg[4:], req)

	ch := make(chan agentResponse)
	a.lock.Lock()
	_, err := a.conn.Write(msg)
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

func (a *sshAgentClient) Add(key agent.AddedKey) error {
	return a.client.Add(key)
}

func (a *sshAgentClient) Remove(key ssh.PublicKey) error {
	return a.client.Remove(key)
}

func (a *sshAgentClient) RemoveAll() error {
	return a.client.RemoveAll()
}

func (a *sshAgentClient) Lock(passphrase []byte) error {
	return a.client.Lock(passphrase)
}

func (a *sshAgentClient) Unlock(passphrase []byte) error {
	return a.client.Unlock(passphrase)
}

func (a *sshAgentClient) Signers() ([]ssh.Signer, error) {
	keys, err := a.client.List()
	if err != nil {
		return nil, err
	}

	ret := make([]ssh.Signer, len(keys))
	for i, k := range keys {
		ret[i] = &sshAgentSigner{a, k}
	}

	return ret, nil
}

func (a *sshAgentClient) SignWithFlags(key ssh.PublicKey, data []byte, flags agent.SignatureFlags) (*ssh.Signature, error) {
	return a.client.SignWithFlags(key, data, flags)
}

func (a *sshAgentClient) Extension(extensionType string, contents []byte) ([]byte, error) {
	return a.client.Extension(extensionType, contents)
}

type sshAgentSigner struct {
	agent *sshAgentClient
	key   ssh.PublicKey
}

func (s *sshAgentSigner) PublicKey() ssh.PublicKey {
	return s.key
}

func (s *sshAgentSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return s.agent.Sign(s.key, data)
}
