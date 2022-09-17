package utils

import (
	"golang.org/x/crypto/ssh"
	"net"
	cfg "press-test/config"
	"time"
)

var (
	sshClient *ssh.Client
)

type IP string

type Connection interface {
	SSHConnect() (*ssh.Client, error)
}

func (i *IP) SSHConnect() (*ssh.Client, error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(cfg.SSHPASSWD))
	clientConfig := &ssh.ClientConfig{
		User:    cfg.SSHUSER,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := string(*i) + ":" + cfg.SSHPORT
	sshClient, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return sshClient, err
	}
	return sshClient, nil
}
