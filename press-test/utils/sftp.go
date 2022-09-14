package utils

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	cfg "press-test/config"
	"time"
)

var (
	sftpClient *sftp.Client
)

func SftpConnect(host string) (*sftp.Client, error){
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
	addr := host + ":" + cfg.SSHPORT
	sshClient, err1 := ssh.Dial("tcp", addr, clientConfig)
	if err1 != nil {
		return sftpClient, err
	}
	//create client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return sftpClient, nil
	}
	return sftpClient, nil
}



