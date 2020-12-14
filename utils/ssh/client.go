package ssh

import (
	"github.houston.softwaregrp.net/Hercules/gotools/utils"
	"github.houston.softwaregrp.net/Hercules/gotools/utils/logging"
	"fmt"
	"github.com/melbahja/goph"
	"os"
	"strings"
)

const maximumAttemptsNumber = 3

func ConnectInteractive(username string, hostname string) *goph.Client {
	logging.LogInfof("\nConnecting to %s...\n", hostname)
	for i := 0; i < maximumAttemptsNumber; i++ {
		if username == "" {
			username = utils.PromptInputFailFast(fmt.Sprintf("Enter user for '%s'", hostname))
		}

		password := utils.PromptPasswordFailFast(fmt.Sprintf("Enter password for '%s' user at '%s'", username, hostname))

		client, err := goph.New(username, hostname, goph.Password(password))
		if err != nil {
			logging.LogErrorf("Cannot connect to %s with provided credentials. Please, try again and make sure they are correct\n", hostname)
			logging.LogError(err)
			continue
		}

		return client
	}

	logging.LogErrorf("Cannot connect to %s with provided credentials. Maximum number of attempts was exceeded: %d", hostname, maximumAttemptsNumber)
	os.Exit(1)
	return nil
}

func Connect(username string, hostname string) (*goph.Client, error) {
	if !utils.FileExists(PrivateKey) {
		logging.LogErrorf("Private key %s is not found. Please, check it was previously generated and has the correct permissions\n", PrivateKey)
		os.Exit(1)
	}
	auth, _ := goph.Key(PrivateKey, "")
	return goph.New(username, hostname, auth)
}

func GetClient(username string, hostname string) *goph.Client {
	var client *goph.Client
	if !strings.EqualFold(hostname, utils.GetCurrentFQDN()) {
		client, _ = Connect(username, hostname)
	}
	return client
}

func SafeCloseClient(client *goph.Client) {
	if client != nil {
		err := client.Close()
		if err != nil {
			logging.LogError(err)
		}
	}
}
