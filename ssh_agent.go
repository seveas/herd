package katyusha

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
}

func (a agentSigner) PublicKey() ssh.PublicKey {
	return a.key
}

func (a agentSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return a.agent.Sign(a.key, data)
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

func sshAgentKeys(path string) ([]ssh.Signer, error) {
	if !agentTried {
		initLock.Lock()
		defer initLock.Unlock()
		if agentTried {
			return agentKeys, agentError
		}
		sock, err := agentConnection()

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
		for i, key := range keys {
			agentKeys[i] = agentSigner{key: key, agent: globalAgent}
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
