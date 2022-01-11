package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	user                string
	password            string
	host                string
	port                int
	backupDirectoryPath string
}

func (c *Connection) connectToRemoteHost() (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(c.password))

	clientConfig = &ssh.ClientConfig{
		User:            c.user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr = fmt.Sprintf("%s:%d", c.host, c.port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func (c *Connection) sendFile(fileToSendPath string) error {
	var (
		err        error
		sftpClient *sftp.Client
	)

	sftpClient, err = c.connectToRemoteHost()
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(fileToSendPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(fileToSendPath)
	dstFile, err := sftpClient.Create(path.Join(c.backupDirectoryPath, remoteFileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	fmt.Println("File sent to server succesfully")
	return nil
}

func (c *Connection) getFile(directoryToCopyTo, fileName string) error {

	var (
		err        error
		sftpClient *sftp.Client
	)

	sftpClient, err = c.connectToRemoteHost()
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(c.backupDirectoryPath + fileName)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var localFileName = path.Base(c.backupDirectoryPath)
	dstFile, err := os.Create(path.Join(directoryToCopyTo, localFileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return err
	}

	fmt.Println("File succesfully copied to ", directoryToCopyTo)

	return nil
}
