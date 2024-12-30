package ssh

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	sshagent "golang.org/x/crypto/ssh/agent"
)

type agentPool struct {
	lock          sync.Mutex
	agents        []sshagent.ExtendedAgent
	signers       []ssh.Signer
	signersByPath map[string][]ssh.Signer
	current       int
}

func newAgentPool(agentCount int, agentTimeout time.Duration) (*agentPool, error) {
	sock, err := agentConnection()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to SSH agent: %s", err)
	}
	a := sshagent.NewClient(sock)
	signers, err := a.Signers()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve keys from SSH agent: %s", err)
	}
	if len(signers) == 0 {
		return nil, fmt.Errorf("No keys found in ssh agent")
	}
	sock, _ = agentConnection()
	pa := newAgent(a, sock)

	pool := &agentPool{
		agents:        make([]sshagent.ExtendedAgent, agentCount),
		signers:       make([]ssh.Signer, len(signers)),
		signersByPath: make(map[string][]ssh.Signer),
	}
	for i, s := range signers {
		pool.signers[i] = &signer{pool, s.PublicKey()}
	}

	if pa.functional(signers[0].PublicKey(), agentTimeout) {
		pool.agents[0] = pa
		for i := 1; i < agentCount; i++ {
			sock, _ = agentConnection()
			pool.agents[i] = newAgent(a, sock)
		}
	} else {
		logrus.Warn("SSH agent is too old, falling back to non-pipelined requests")
		pool.agents[0] = a
		for i := 1; i < agentCount; i++ {
			sock, _ = agentConnection()
			pool.agents[i] = sshagent.NewClient(sock)
		}
	}
	return pool, nil
}

func (ap *agentPool) nextAgent() sshagent.ExtendedAgent {
	ap.lock.Lock()
	defer func() {
		ap.current = (ap.current + 1) % len(ap.agents)
		ap.lock.Unlock()
	}()
	return ap.agents[ap.current]
}

func (ap *agentPool) List() ([]*sshagent.Key, error) {
	return ap.nextAgent().List()
}

func (ap *agentPool) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	return ap.nextAgent().Sign(key, data)
}

func (ap *agentPool) SignWithFlags(key ssh.PublicKey, data []byte, flags sshagent.SignatureFlags) (*ssh.Signature, error) {
	return ap.nextAgent().SignWithFlags(key, data, flags)
}

func (ap *agentPool) Extension(extensionType string, contents []byte) ([]byte, error) {
	return ap.nextAgent().Extension(extensionType, contents)
}

func (ap *agentPool) Add(key sshagent.AddedKey) error {
	return ap.nextAgent().Add(key)
}

func (ap *agentPool) Remove(key ssh.PublicKey) error {
	return ap.nextAgent().Remove(key)
}

func (ap *agentPool) RemoveAll() error {
	return ap.nextAgent().RemoveAll()
}

func (ap *agentPool) Lock(passphrase []byte) error {
	return ap.nextAgent().Lock(passphrase)
}

func (ap *agentPool) Unlock(passphrase []byte) error {
	return ap.nextAgent().Unlock(passphrase)
}

func (ap *agentPool) Signers() ([]ssh.Signer, error) {
	return ap.signers, nil
}

func (ap *agentPool) SignersForPathCallback(path string) func() ([]ssh.Signer, error) {
	return func() ([]ssh.Signer, error) {
		signers := ap.SignersForPath(path)
		if len(signers) == 0 {
			return nil, fmt.Errorf("SSH key %s was not found in the SSH agent", path)
		}
		return signers, nil
	}
}

func (ap *agentPool) SignersForPath(path string) []ssh.Signer {
	if path == "" {
		return ap.signers
	}
	if k, ok := ap.signersByPath[path]; ok {
		return k
	}
	ap.lock.Lock()
	defer ap.lock.Unlock()
	for _, signer := range ap.signers {
		if signer.PublicKey().(*sshagent.Key).Comment == path {
			ap.signersByPath[path] = []ssh.Signer{signer}
			return []ssh.Signer{signer}
		}
	}

	// If we didn't find the key, try again by parsing the public key and matching by key data
	ap.signersByPath[path] = []ssh.Signer{}
	data, err := os.ReadFile(path + ".pub")
	if err != nil {
		return []ssh.Signer{}
	}
	key, _, _, _, err := ssh.ParseAuthorizedKey(data) //nolint:dogsled // Can't help it that we don't need the rest
	if err != nil {
		return []ssh.Signer{}
	}
	mkey := key.Marshal()
	for _, signer := range ap.signers {
		if bytes.Equal(signer.PublicKey().Marshal(), mkey) {
			ap.signersByPath[path] = []ssh.Signer{signer}
			return []ssh.Signer{signer}
		}
	}
	return []ssh.Signer{}
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

var _ sshagent.ExtendedAgent = &agentPool{}

type signer struct {
	agent sshagent.ExtendedAgent
	key   ssh.PublicKey
}

func (s *signer) PublicKey() ssh.PublicKey {
	return s.key
}

func (s *signer) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return s.agent.Sign(s.key, data)
}

var _ ssh.Signer = &signer{}
