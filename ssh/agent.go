package ssh

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	sshagent "golang.org/x/crypto/ssh/agent"
)

type agent struct {
	sshAgent   sshagent.ExtendedAgent
	connection io.ReadWriter
	waiters    []chan agentResponse
	lock       sync.Mutex
}

func newAgent(sshAgent sshagent.ExtendedAgent, sock io.ReadWriter) *agent {
	a := &agent{
		sshAgent:   sshAgent,
		connection: sock,
		waiters:    make([]chan agentResponse, 0, 64),
	}
	go a.readLoop()

	return a
}

type agentResponse struct {
	data []byte
	err  error
}

// Test the ssh agent by sending a bunch of requests in a pipelined way. If
// they are not answered within the specified interval (50ms by default), the
// ssh agent is too old and suffers from the bug solved in
// https://github.com/openssh/openssh-portable/pull/183
func (a *agent) functional(key ssh.PublicKey, timeout time.Duration) bool {
	tests := 10
	c := make(chan bool)
	t := time.NewTicker(timeout)
	defer t.Stop()
	for range tests {
		go func() { _, err := a.Sign(key, []byte("Test")); c <- (err == nil) }()
	}
	for range tests {
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
	if _, err := io.ReadFull(a.connection, respSizeBuf[:]); err != nil {
		return nil, err
	}
	respSize := binary.BigEndian.Uint32(respSizeBuf[:])
	buf := make([]byte, respSize)
	if _, err := io.ReadFull(a.connection, buf); err != nil {
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
	return a.SignWithFlags(key, data, 0)
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
	return a.sshAgent.Signers()
}

func (a *agent) SignWithFlags(key ssh.PublicKey, data []byte, flags sshagent.SignatureFlags) (*ssh.Signature, error) {
	req := ssh.Marshal(agentSignRequest{Key: key.Marshal(), Data: data, Flags: uint32(flags)})
	if len(req) > math.MaxUint32 {
		return nil, errors.New("ssh agent request too large")
	}
	msg := make([]byte, 4+len(req))
	binary.BigEndian.PutUint32(msg, uint32(len(req))) // nolint:gosec // Overflow is checked above
	copy(msg[4:], req)

	ch := make(chan agentResponse)
	a.lock.Lock()
	_, err := a.connection.Write(msg)
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

func (a *agent) Extension(extensionType string, contents []byte) ([]byte, error) {
	return a.sshAgent.Extension(extensionType, contents)
}

var _ sshagent.ExtendedAgent = &agent{}
