package app

import (
	"context"
	"github.com/sirupsen/logrus"
	kubernetes "k8s-from-secrets-vault/kubernetes"
	vault "k8s-from-secrets-vault/vault"
	"os"
)

type Command struct {
	Address        string
	AuthToken      string
	VaultNamespace string
	EngineName     string
	SecretPath     string

	Base64Kubeconfig string
	Namespace        string

	LoadAsConfigMap bool

	kubernetesClient kubernetes.KubernetesClient
}

func SetupCommand(args []string) Command {
	return Command{
		Address:          args[0],
		AuthToken:        args[1],
		VaultNamespace:   args[2],
		EngineName:       args[3],
		SecretPath:       args[4],
		Base64Kubeconfig: args[5],
		Namespace:        args[6],
		LoadAsConfigMap:  args[7] == "true",
	}
}

func SetupCommandWithKubernetesClient(args []string, kubernetesClient kubernetes.KubernetesClient) Command {
	command := SetupCommand(args)
	command.kubernetesClient = kubernetesClient
	return command
}

func (command Command) Execute() error {
	if command.LoadAsConfigMap {
		return command.loadAndApplyConfigMap()
	} else {
		return command.loadAndApplySecrets()
	}
}

func (command Command) vaultParameters() vault.VaultConfig {
	return vault.VaultConfig{
		Address:    command.Address,
		AuthToken:  command.AuthToken,
		Namespace:  command.VaultNamespace,
		EngineName: command.EngineName,
		SecretPath: command.SecretPath,
	}
}

func (command Command) kubeParameters() kubernetes.KubernetesParameters {
	return kubernetes.KubernetesParameters{
		Base64Kubeconfig: command.Base64Kubeconfig,
		Namespace:        command.Namespace,
	}
}

func (command Command) loadAndApplySecrets() error {
	log := setupLogger()

	data, err := vault.LoadSecretData(command.vaultParameters(), log)

	if err != nil {
		return err
	}

	kubernetesConfig, err := kubernetes.CreateConfig(command.kubeParameters(), log)
	if err != nil {
		return err
	}

	kubernetesClient, err := kubernetes.CreateClient(kubernetesConfig, log)
	if command.kubernetesClient != nil {
		kubernetesClient = command.kubernetesClient
	}
	if err != nil {
		return err
	}

	err = kubernetesClient.ApplySecret(context.TODO(), "app-secret", data, log)
	if err != nil {
		return err
	}

	return nil
}

func (command Command) loadAndApplyConfigMap() error {
	log := setupLogger()

	data, err := vault.LoadSecretData(command.vaultParameters(), log)

	if err != nil {
		return err
	}

	kubernetesConfig, err := kubernetes.CreateConfig(command.kubeParameters(), log)
	if err != nil {
		return err
	}

	kubernetesClient, err := kubernetes.CreateClient(kubernetesConfig, log)
	if command.kubernetesClient != nil {
		kubernetesClient = command.kubernetesClient
	}
	if err != nil {
		return err
	}

	err = kubernetesClient.ApplyConfigMap(context.TODO(), "app-config", data, log)
	if err != nil {
		return err
	}

	return nil
}

func setupLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{}
	return log
}
