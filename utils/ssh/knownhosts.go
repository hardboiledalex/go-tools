package ssh

import (
	"github.houston.softwaregrp.net/Hercules/gotools/utils"
	"github.houston.softwaregrp.net/Hercules/gotools/utils/logging"
	"encoding/base64"
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strings"
)

const (
	defaultSSHPort = 22
)

// Adds remote host to the local known_hosts file
func AddHostToLocalKnownHosts(hostname string, localKnownHostsPath string) {
	sshConfig := &ssh.ClientConfig{
		HostKeyCallback: addLocalKnownHostCallback(localKnownHostsPath),
	}
	ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, defaultSSHPort), sshConfig)
}

func addLocalKnownHostCallback(localKnownHostsPath string) ssh.HostKeyCallback {
	return func(dialAddr string, addr net.Addr, publicKey ssh.PublicKey) error {
		knownHostsFile, err := os.OpenFile(localKnownHostsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			logging.LogError(err)
		}

		knownHostsFileString := utils.ReadFileToString(localKnownHostsPath)
		newKnownHostLine := fmt.Sprintf("%s %s %s", strings.Split(dialAddr, ":")[0], publicKey.Type(), base64.StdEncoding.EncodeToString(publicKey.Marshal()))
		if !strings.Contains(knownHostsFileString, newKnownHostLine) {
			if _, err := knownHostsFile.WriteString("\n" + newKnownHostLine + "\n"); err != nil {
				logging.LogError(err)
			}
		}

		defer knownHostsFile.Close()

		return nil
	}
}

// Adds remote host to the remote known_hosts file
func AddHostToRemoteKnownHosts(client *goph.Client, hostname string, remoteKnownHostsPath string) {
	sshConfig := &ssh.ClientConfig{
		HostKeyCallback: addRemoteKnownHostCallback(client, remoteKnownHostsPath),
	}
	ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, defaultSSHPort), sshConfig)
}

func addRemoteKnownHostCallback(client *goph.Client, remoteKnownHostsPath string) ssh.HostKeyCallback {
	return func(dialAddr string, addr net.Addr, publicKey ssh.PublicKey) error {
		remoteKnownHostsFile := DownloadTextFileToMemory(client, remoteKnownHostsPath)
		knownHostsFileString := strings.Join(remoteKnownHostsFile.Strings[:], "\n")
		newKnownHostLine := fmt.Sprintf("%s %s %s", strings.Split(dialAddr, ":")[0], publicKey.Type(), base64.StdEncoding.EncodeToString(publicKey.Marshal()))
		if !strings.Contains(knownHostsFileString, newKnownHostLine) {
			knownHostsFileString += "\n" + newKnownHostLine + "\n"
			err := UploadTextFileFromMemory(client, remoteKnownHostsPath, &TextFile{
				Strings: strings.Split(knownHostsFileString, "\n"),
				Mode:    remoteKnownHostsFile.Mode,
			})
			return err
		}

		return nil
	}
}
