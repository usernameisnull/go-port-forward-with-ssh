package main

import (
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"port_forward_mine/config"
	"port_forward_mine/logs"
	"time"
)

var ConfigInfo config.Config

func init() {
	ConfigInfo = config.ConfigParse()
	logs.SetLog(ConfigInfo.LogPath, 10, 10, 10)
}

// Get default location of a private key
func privateKeyPath() string {
	return ConfigInfo.Pem
}

// Get private key for ssh authentication
func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, _ := ioutil.ReadFile(keyPath)
	return ssh.ParsePrivateKeyWithPassphrase(buff, []byte(ConfigInfo.PemPass))
}

// Get ssh client config for our connection
// SSH config will use 2 authentication strategies: by key and by password
func makeSshConfig(user string) (*ssh.ClientConfig, error) {
	key, err := parsePrivateKey(privateKeyPath())
	if err != nil {
		return nil, err
	}

	sshConfig := ssh.ClientConfig{
		User: user,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//fmt.Println(hostname, remote, key)
			return nil
		},
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: time.Second * 600,
	}

	return &sshConfig, nil
}

// Handle local client connections and tunnel data to the remote serverq
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, cfg *ssh.ClientConfig) {

	// Establish connection with SSH server
	conn, err := ssh.Dial("tcp", ConfigInfo.JumpServer, cfg)
	if err != nil {
		logs.Error.Println(err)
	}
	defer conn.Close()

	// Establish connection with remote server
	remote, err := conn.Dial("tcp", ConfigInfo.Remote)
	if err != nil {
		logs.Error.Println(err)
	}
	chDone := make(chan bool)
	defer client.Close()
	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			logs.Error.Println("error while copy remote->local:", err)
		}
		//fmt.Println("remote->local copy 完成")
		chDone <- true
	}()
	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			logs.Error.Println(err)
		}
		//fmt.Println("local->remote copy 完成")
		chDone <- true
	}()
	<-chDone
}

func main() {
	logs.Info.Printf("开始,在 %s<------>%s 之间进行转发", ConfigInfo.Remote, ConfigInfo.Local)
	// Build SSH client configuration
	cfg, err := makeSshConfig(ConfigInfo.JumpServerUser)
	if err != nil {
		logs.Error.Println(err)
	}

	// Start local server to forward traffic to remote connection
	local, err := net.Listen("tcp", ConfigInfo.Local)
	if err != nil {
		logs.Error.Println(err)
	}
	defer local.Close()
	if local == nil {
		logs.Error.Printf("dial %s get nil", ConfigInfo.Local)
		return
	}
	// Handle incoming connections
	for {
		client, err := local.Accept()
		if err != nil {
			logs.Error.Println(err)
		}
		if client == nil {
			logs.Error.Println("local.Accept() get nil")
			return
		}
		go handleClient(client, cfg)
	}
}
