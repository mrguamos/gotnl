package gotnl

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"log"

	"golang.org/x/crypto/ssh"
)

type Config struct {
	BastionHost string
	BastionPort string
	BastionUser string
	TargetHost  string
	TargetPort  string
	LocalPort   string
	SSHKey      string
	Passphrase  string
}

func Tunnel(cfg Config) (*ssh.Client, net.Listener, error) {

	authMethod, err := sshAgentAuth(cfg)
	if err != nil {
		return nil, nil, err
	}

	config := &ssh.ClientConfig{
		User: cfg.BastionUser,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
	}

	bastionAddr := fmt.Sprintf("%s:%s", cfg.BastionHost, cfg.BastionPort)
	log.Println("connecting to bastion host")
	bastionConn, err := ssh.Dial("tcp", bastionAddr, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial bastion host: %v", err)
	}

	localAddr := fmt.Sprintf("localhost:%s", cfg.LocalPort)
	targetAddr := fmt.Sprintf("%s:%s", cfg.TargetHost, cfg.TargetPort)
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create listener on bastion host: %v", err)
	}

	log.Printf("SSH tunnel established: %s -> %s", localAddr, targetAddr)

	go func() {
		for {
			localConn, err := listener.Accept()
			if err != nil {
				fmt.Println("Failed to accept connection:", err)
				listener.Close()
				bastionConn.Close()
				break
			}

			go handleLocalConnection(bastionConn, localConn, targetAddr)
		}
	}()

	return bastionConn, listener, nil
}

func publicKeyFile(file string, cfg Config) (ssh.AuthMethod, error) {
	key, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	passPhrase := cfg.Passphrase
	var signer ssh.Signer

	if len(passPhrase) > 0 {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passPhrase))
	} else {
		signer, err = ssh.ParsePrivateKey(key)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

func handleLocalConnection(client *ssh.Client, localConn net.Conn, targetAddr string) {
	targetConn, err := client.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target host: %v", err)
		return
	}

	go copyWithCheck(localConn, targetConn)
	go copyWithCheck(targetConn, localConn)
}

func copyWithCheck(dst net.Conn, src net.Conn) {
	_, err := io.Copy(dst, src)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Println("Connection closed by client")
		} else {
			log.Printf("Failed to copy data: %v", err)
		}
	}
}

func sshAgentAuth(cfg Config) (ssh.AuthMethod, error) {
	sshKey := cfg.SSHKey

	if len(sshKey) > 0 {
		return publicKeyFile(sshKey, cfg)
	}

	return nil, errors.New("missing SSH_KEY, please configure it in the settings")
}
