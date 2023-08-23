package main

import (
	"k8s-from-secrets-vault/app"
	"os"
)

func main() {
	command := app.SetupCommand()
	err := command.Execute()
	if err != nil {
		os.Exit(1)
		return
	}
	os.Exit(0)
}
