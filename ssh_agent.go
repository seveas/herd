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

type agentSigner struct {
	agent ssh_agent.Agent
	key   ssh.PublicKey
	lock  *sync.Mutex
}

func (a agentSigner) PublicKey() ssh.PublicKey {
	return a.key
}

func (a agentSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	a.lock.Lock()
	signature, err := a.agent.Sign(a.key, data)
	a.lock.Unlock()
	return signature, err
}

func sshAgentKeys(path string) ([]ssh.Signer, error) {
	if !agentTried {
		initLock.Lock()
		defer initLock.Unlock()
		if agentTried {
			return agentKeys, agentError
		}
		sockPath, ok := os.LookupEnv("SSH_AUTH_SOCK")
		if !ok {
			agentError = fmt.Errorf("No ssh agent found in environment, make sure your ssh agent is running")
			if _, ok = os.LookupEnv("SSH_CONNECTION"); ok {
				agentError = fmt.Errorf("No ssh agent found in environment, make sure your ssh agent is running and forwarded")
			}
			return agentKeys, agentError
		}

		sock, err := net.Dial("unix", sockPath)
		if err != nil {
			agentError = fmt.Errorf("Unable to connect to SSH agent: %s", err)
			return agentKeys, agentError
		}

		globalAgent = ssh_agent.NewClient(sock)

		keys, err := globalAgent.List()
		if err != nil {
			agentError = fmt.Errorf("Unable to retrieve keys from SSH agent: %s", err)
			return agentKeys, agentError
		}

		agentKeys = make([]ssh.Signer, len(keys))
		var lock sync.Mutex
		for i, key := range keys {
			agentKeys[i] = agentSigner{key: key, agent: globalAgent, lock: &lock}
		}
		agentTried = true
	}
	if path != "" {
		for _, k := range agentKeys {
			if k.(agentSigner).key.(*ssh_agent.Key).Comment == path {
				return []ssh.Signer{k}, agentError
			}
		}
		return []ssh.Signer{}, fmt.Errorf("Key %s not found", path)
	}
	return agentKeys, agentError
}
