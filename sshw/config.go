package sshw

import (
	"golang.org/x/crypto/ssh"
)

type Node struct {
	Name       string  `json:"name"`
	Host       string  `json:"host"`
	User       string  `json:"user"`
	Port       int     `json:"port"`
	KeyPath    string  `json:"keypath"`
	Passphrase string  `json:"passphrase"`
	Password   string  `json:"password"`
	Children   []*Node `json:"children"`
}

func (n *Node) String() string {
	return n.Name
}

func (n *Node) user() string {
	if n.User == "" {
		return "root"
	}
	return n.User
}

func (n *Node) port() int {
	if n.Port <= 0 {
		return 22
	}
	return n.Port
}

func (n *Node) password() ssh.AuthMethod {
	if n.Password == "" {
		return nil
	}
	return ssh.Password(n.Password)
}

var (
	config []*Node
)

func GetConfig() []*Node {
	return config
}
