package sshw

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	DEFAULT_CHUNK_SIZE        = 65536
	MIN_CHUNKS                = 10
	THROUGHPUT_SLEEP_INTERVAL = 100
	MIN_THROUGHPUT            = DEFAULT_CHUNK_SIZE * MIN_CHUNKS * (1000 / THROUGHPUT_SLEEP_INTERVAL)
)

var (
	maxThroughputChan  = make(chan bool, MIN_CHUNKS)
	maxThroughput      uint64
	maxThroughputMutex sync.Mutex
)

type Client interface {
	Login()
	Cmd(string) ([]byte, error)
	UploadFile(sourceFile, target string) (stdout, stderr string, err error)
}

type defaultClient struct {
	clientConfig *ssh.ClientConfig
	node         *Node
}

func NewClient(node *Node) Client {
	u, err := user.Current()
	if err != nil {
		l.Error(err)
		return nil
	}

	var authMethods []ssh.AuthMethod

	var pemBytes []byte
	if node.KeyPath == "" {
		pemBytes, err = ioutil.ReadFile(path.Join(u.HomeDir, ".ssh/id_rsa"))
	} else {
		pemBytes, err = ioutil.ReadFile(node.KeyPath)
	}
	if err != nil {
		l.Error(err)
	} else {
		var signer ssh.Signer
		if node.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(node.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		}
		if err != nil {
			l.Error(err)
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	password := node.password()

	if password != nil {
		interactive := func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			answers = make([]string, len(questions))
			for n := range questions {
				answers[n] = node.Password
			}

			return answers, nil
		}
		authMethods = append(authMethods, ssh.KeyboardInteractive(interactive))
		authMethods = append(authMethods, password)
	}

	config := &ssh.ClientConfig{
		User:            node.user(),
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	config.SetDefaults()
	config.Ciphers = append(config.Ciphers, "aes128-cbc", "3des-cbc", "blowfish-cbc", "cast128-cbc", "aes192-cbc", "aes256-cbc")

	return &defaultClient{
		clientConfig: config,
		node:         node,
	}
}

func (c *defaultClient) UploadFile(sourceFile, target string) (stdout, stderr string, err error) {
	log.Println("cp file:" + sourceFile + " to " + target)
	host := c.node.Host
	port := c.node.port()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), c.clientConfig)
	if err != nil {
		l.Error(err)
		panic(err)
	}
	defer client.Close()

	currentSession, err := client.NewSession()
	if err != nil {
		l.Error(err)
		panic(err)
	}
	defer currentSession.Close()

	l.Infof("connect server ssh -p %d %s@%s version: %s\n", port, c.node.user(), host, string(client.ServerVersion()))

	f, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	cmd := "cat >'" + strings.Replace(target, "'", "'\\''", -1) + "'"
	stdinPipe, err := currentSession.StdinPipe()
	if err != nil {
		panic(err)
	}
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	currentSession.Stdout = &stdoutBuf
	currentSession.Stderr = &stderrBuf

	err = currentSession.Start(cmd)
	if err != nil {
		panic(err)
	}

	for start, max := 0, len(data); start < max; start += DEFAULT_CHUNK_SIZE {
		// <-maxThroughputChan
		end := start + DEFAULT_CHUNK_SIZE
		if end > max {
			end = max
		}

		_, err = stdinPipe.Write(data[start:end])
		if err != nil {
			panic(err)
		}
	}

	err = stdinPipe.Close()
	if err != nil {
		panic(err)
	}
	err = currentSession.Wait()
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	return
}

func (c *defaultClient) Cmd(cmd string) ([]byte, error) {
	host := c.node.Host
	port := c.node.port()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), c.clientConfig)
	if err != nil {
		l.Error(err)
		panic(err)
	}
	defer client.Close()

	l.Infof("connect server ssh -p %d %s@%s version: %s\n", port, c.node.user(), host, string(client.ServerVersion()))

	session, err := client.NewSession()
	if err != nil {
		l.Error(err)
		panic(err)
	}
	defer session.Close()

	output, error := session.CombinedOutput(cmd)
	return output, error
}

func (c *defaultClient) Login() {
	host := c.node.Host
	port := c.node.port()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), c.clientConfig)
	if err != nil {
		l.Error(err)
		return
	}
	defer client.Close()

	l.Infof("connect server ssh -p %d %s@%s version: %s\n", port, c.node.user(), host, string(client.ServerVersion()))

	session, err := client.NewSession()
	if err != nil {
		l.Error(err)
		return
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		l.Error(err)
		return
	}
	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		l.Error(err)
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", h, w, modes)
	if err != nil {
		l.Error(err)
		return
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Shell()
	if err != nil {
		l.Error(err)
		return
	}

	session.Wait()
}
