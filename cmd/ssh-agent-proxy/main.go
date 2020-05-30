package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	// Make sure we don't try to proxy ourselves
	if _, ok := os.LookupEnv("SSH_AGENT_PROXY"); ok {
		os.Exit(0)
	}
	os.Setenv("SSH_AGENT_PROXY", "1")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "ssh-agent-proxy",
	Long:    "SSH agent proxy",
	Version: "0.1",
	Run:     runProxy,
}

func init() {
	rootCmd.Flags().StringP("address", "a", "", "Bind to this socket")
	rootCmd.Flags().BoolP("foreground", "f", false, "Stay in the foreground")
	rootCmd.Flags().BoolP("debug", "d", false, "Stay in the foreground and print debug information")
	rootCmd.Flags().SetInterspersed(false)
}

func runProxy(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	addr, ok := os.LookupEnv("SSH_AUTH_SOCK")
	if !ok {
		logrus.Fatalf("No upstream SSH agent could be found")
	}

	path, _ := cmd.Flags().GetString("address")
	if path == "" {
		dir, err := ioutil.TempDir("", "ssh-agent-proxy-*")
		if err != nil {
			logrus.Fatalf("Could not create temporary directory: %s", err)
		}
		path = filepath.Join(dir, "agent.sock")
	}

	fmt.Printf("SSH_AUTH_SOCK='%s'; export SSH_AUTH_SOCK\n", path)
	fmt.Println("SSH_AGENT_PROXY=1; export SSH_AGENT_PROXY")
	fmt.Println("KATYUSHA_FAST_SSH_AGENT=1; export KATYUSHA_FAST_SSH_AGENT")

	if fg, _ := cmd.Flags().GetBool("foreground"); !debug && !fg {
		null, err := os.Open(os.DevNull)
		if err != nil {
			panic(fmt.Sprintf("Unable to open %s", os.DevNull))
		}
		attrs := syscall.ProcAttr{
			Files: []uintptr{
				null.Fd(),
				null.Fd(),
				null.Fd(),
			},
			Env: os.Environ(),
		}
		exe, err := exec.LookPath(os.Args[0])
		if err != nil {
			logrus.Fatalf("Unable to launch into the background: %s", err)
		}
		if _, err := syscall.ForkExec(os.Args[0], []string{exe, "--address", path, "--foreground"}, &attrs); err != nil {
			fmt.Println(os.Args)
			logrus.Fatalf("Unable to launch into the background: %s", err)
		}
		if len(args) != 0 {
			os.Setenv("SSH_AUTH_SOCK", path)
			os.Setenv("KATYUSHA_FAST_SSH_AGENT=1", path)
			exe, err := exec.LookPath(args[0])
			if err != nil {
				exe = args[0]
			}
			if err := syscall.Exec(exe, args, os.Environ()); err != nil {
				logrus.Fatalf("Unable to launch %s: %s", exe, err)
			}
		}
		os.Exit(0)
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		logrus.Fatalf("Unable to listen to unix socket on %s: %s", path, err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			logrus.Fatalf("Unable to accept new client: %s", err)
		}
		logrus.Infof("New client connected")
		agent, err := net.Dial("unix", addr)
		if err != nil {
			logrus.Fatalf("Could not connect to upstream ssh agent: %s", err)
		}
		go clientLoop(client, agent)
	}
}

func clientLoop(client, agent io.ReadWriteCloser) {
	defer client.Close()
	defer agent.Close()

	messages := make(chan []byte, 10240)
	go func() {
		defer close(messages)
		for {
			// Read single message from client
			msg, err := readSingleMessage(client)
			if err != nil {
				logrus.Infof("Error reading from client: %s", err)
				break
			}
			logrus.Debugf("Received message from client. Type %d, length %d", int(msg[4]), len(msg)-4)
			messages <- msg
		}
	}()

	for {
		// Forward to agent
		msg, ok := <-messages
		if !ok {
			logrus.Infof("Message channel closed")
			break
		}
		if _, err := agent.Write(msg); err != nil {
			logrus.Infof("Error writing to agent: %s", err)
			break
		}

		// Read response
		msg, err := readSingleMessage(agent)
		if err != nil {
			logrus.Infof("Error reading from agent: %s", err)
			break
		}
		logrus.Debugf("Received message from agent. Type %d, length %d", int(msg[4]), len(msg)-4)

		// Forward to client
		if _, err := client.Write(msg); err != nil {
			logrus.Infof("Error writing to client: %s", err)
			break
		}
	}
}

func readSingleMessage(sock io.Reader) ([]byte, error) {
	lbuf := make([]byte, 4)
	if _, err := io.ReadFull(sock, lbuf); err != nil {
		return nil, err
	}
	mbuf := make([]byte, binary.BigEndian.Uint32(lbuf))
	if _, err := io.ReadFull(sock, mbuf); err != nil {
		return nil, err
	}
	return append(lbuf, mbuf...), nil
}
