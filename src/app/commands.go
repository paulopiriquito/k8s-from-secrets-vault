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

func (command Command) LoadAndApplySecrets() error {
	log := setupLogger()

	data, err := vault.LoadSecretData(command.vaultParameters(), log)

	if err != nil {
		return err
	}

	kubernetesConfig, err := kubernetes.NewKubernetesConfig(command.kubeParameters())
	if err != nil {
		return err
	}

	kubernetesClient, err := kubernetes.NewKubernetesClient(kubernetesConfig)
	if err != nil {
		return err
	}

	err = kubernetesClient.ApplySecret(context.TODO(), "app-secret", data)
	if err != nil {
		return err
	}

	return nil
}

func (command Command) LoadAndApplyConfigMap() interface{} {
	log := setupLogger()

	data, err := vault.LoadSecretData(command.vaultParameters(), log)

	if err != nil {
		return err
	}

	kubernetesConfig, err := kubernetes.NewKubernetesConfig(command.kubeParameters())
	if err != nil {
		return err
	}

	kubernetesClient, err := kubernetes.NewKubernetesClient(kubernetesConfig)
	if err != nil {
		return err
	}

	err = kubernetesClient.ApplyConfigMap(context.TODO(), "app-config", data)
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
