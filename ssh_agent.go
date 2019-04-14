package herd

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
	ssh_agent "golang.org/x/crypto/ssh/agent"
)

var globalAgent ssh_agent.Agent
var agentKeys []ssh.Signer
var agentTried bool = false
var agentError error = nil
var initLock sync.Mutex

type AgentSigner struct {
	agent ssh_agent.Agent
	key   ssh.PublicKey
	lock  *sync.Mutex
}

func (a AgentSigner) PublicKey() ssh.PublicKey {
	return a.key
}

func (a AgentSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	a.lock.Lock()
	signature, err := a.agent.Sign(a.key, data)
	a.lock.Unlock()
	return signature, err
}

func SshAgentKeys() ([]ssh.Signer, error) {
	if !agentTried {
		initLock.Lock()
		defer initLock.Unlock()
		if agentTried {
			return agentKeys, agentError
		}
		sockPath, ok := os.LookupEnv("SSH_AUTH_SOCK")
		if !ok {
			agentError = fmt.Errorf("No ssh agent found in environment")
			return agentKeys, agentError
		}

		sock, err := net.Dial("unix", sockPath)
		if err != nil {
			agentError = fmt.Errorf("Unable to connect to SSH agent: %s\n", err)
			return agentKeys, agentError
		}

		globalAgent = ssh_agent.NewClient(sock)

		keys, err := globalAgent.List()
		if err != nil {
			agentError = fmt.Errorf("Unable to retrieve keys from SSH agent: %s\n", err)
			return agentKeys, agentError
		}

		agentKeys = make([]ssh.Signer, len(keys))
		var lock sync.Mutex
		for i, key := range keys {
			agentKeys[i] = AgentSigner{key: key, agent: globalAgent, lock: &lock}
		}
		agentTried = true
	}
	return agentKeys, agentError
}
