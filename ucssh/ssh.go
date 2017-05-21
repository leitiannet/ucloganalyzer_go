package ucssh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

/*
使用Go语言实现远程传输文件
http://www.aspku.com/tech/jiaoben/golang/196410.html
https://github.com/pkg/sftp/tree/master/examples
https://github.com/rapidloop/rtop
*/

type SSHClient struct {
	conn *ssh.Client
}

func NewSSHClient() *SSHClient {
	client := &SSHClient{}
	return client
}

func (this *SSHClient) Dial(host string, port int, user, passwd string) bool {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return false
	}
	this.conn = conn
	return true
}

func (this *SSHClient) Run(command string) (string, error) {
	if this.conn == nil {
		return "", fmt.Errorf("conn nil")
	}
	conn := this.conn
	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var result bytes.Buffer
	session.Stdout = &result
	session.Stderr = &result
	session.Stdin = &result
	err = session.Run(command)
	if err != nil {
		return "", err
	}
	return string(result.Bytes()), nil
}

func (this *SSHClient) Upload(localFile, remoteDir string) error {
	if this.conn == nil {
		return fmt.Errorf("conn nil")
	}
	conn := this.conn

	sftpClient, err := sftp.NewClient(conn, sftp.MaxPacket(6e9))
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	var localBaseFile = path.Base(localFile)
	srcFile, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := sftpClient.Create(path.Join(remoteDir, localBaseFile))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 100*1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	fmt.Println("copy file to remote server finished!")
	return nil
}

func (this *SSHClient) Download(remoteFile, localDir string) error {
	if this.conn == nil {
		return fmt.Errorf("conn nil")
	}
	conn := this.conn

	sftpClient, err := sftp.NewClient(conn, sftp.MaxPacket(6e9))
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	var remoteBaseFile = path.Base(remoteFile)
	srcFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(path.Join(localDir, remoteBaseFile))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	info, _ := srcFile.Stat()
	io.Copy(dstFile, io.LimitReader(srcFile, info.Size()))

	fmt.Println("copy file from remote server finished!")
	return nil
}

func (this *SSHClient) Close() {
	if this.conn != nil {
		this.conn.Close()
	}
	this.conn = nil
}
