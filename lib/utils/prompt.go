package utils

import (
	"github.com/hardboiledalex/go-tools/lib/logging"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"os"
	"strings"
)

func PromptInput(message string) (string, error) {
	output := ""
	prompt := &survey.Input{
		Message: message,
	}
	err := survey.AskOne(prompt, &output)
	return output, err
}

func PromptInputFailFast(message string) string {
	output, err := PromptInput(message)
	if err != nil {
		logging.LogError(err)
		os.Exit(1)
	}
	return output
}

func PromptPassword(message string) (string, error) {
	password := ""
	prompt := &survey.Password{
		Message: message,
	}
	err := survey.AskOne(prompt, &password)
	return password, err
}

func PromptPasswordFailFast(message string) string {
	for {
		password, err := PromptPassword(message)
		if err != nil {
			logging.LogError(err)
			os.Exit(1)
		}
		trimmedPassword := strings.TrimSpace(password)
		if trimmedPassword == "" {
			color.Red.Println("Password cannot be empty. Please, try again.")
			continue
		}
		return trimmedPassword
	}
}

func PromptConfirm(message string) (bool, error) {
	output := false
	prompt := &survey.Confirm{
		Message: message,
	}
	err := survey.AskOne(prompt, &output)
	return output, err
}

func PromptConfirmFailFast(message string) bool {
	output, err := PromptConfirm(message)
	if err != nil {
		logging.LogError(err)
		os.Exit(1)
	}
	return output
}

func PromptSelect(message string, inputOptions []string) (string, error) {
	output := ""
	prompt := &survey.Select{
		Message: message,
		Options: inputOptions,
	}
	err := survey.AskOne(prompt, &output)
	return output, err
}

func PromptSelectFailFast(message string, inputOptions []string) string {
	output, err := PromptSelect(message, inputOptions)
	if err != nil {
		logging.LogError(err)
		os.Exit(1)
	}
	return output
}

func PromptMultiSelect(message string, inputOptions []string) ([]string, error) {
	var outputOptions []string
	prompt := &survey.MultiSelect{
		Message: message,
		Options: inputOptions,
	}
	err := survey.AskOne(prompt, &outputOptions)
	return outputOptions, err
}

func PromptMultiSelectFailFast(message string, inputOptions []string) []string {
	outputOptions, err := PromptMultiSelect(message, inputOptions)
	if err != nil {
		logging.LogError(err)
		os.Exit(1)
	}
	return outputOptions
}
