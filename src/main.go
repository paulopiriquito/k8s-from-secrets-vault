package main

import (
	"k8s-from-secrets-vault/app"
	"os"
)

func main() {
	command, err := app.SetupCommand()
	if err != nil {
		os.Exit(1)
		return
	}

	err = command.Execute()
	if err != nil {
		os.Exit(1)
		return
	}
	os.Exit(0)
}
