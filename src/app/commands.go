package app

import (
	"context"
	"github.com/sirupsen/logrus"
	kubernetes "k8s-from-secrets-vault/kubernetes"
	vault "k8s-from-secrets-vault/vault"
	"os"
)

const (
	VaultAddress      = "VAULT_ADDRESS"
	VaultToken        = "VAULT_TOKEN"
	VaultNamespace    = "VAULT_NAMESPACE"
	VaultEngine       = "VAULT_ENGINE"
	VaultSecretPath   = "VAULT_SECRET_PATH"
	Kubeconfig        = "KUBECONFIG"
	Namespace         = "KUBERNETES_NAMESPACE"
	ApplyAsConfigmap  = "LOAD_AS_CONFIGMAP"
	ObjectNameToApply = "OBJECT_NAME_TO_APPLY"
)

type Command struct {
	Address        string
	AuthToken      string
	VaultNamespace string
	EngineName     string
	SecretPath     string

	Base64Kubeconfig string
	Namespace        string

	LoadAsConfigMap   bool
	ObjectNameToApply string

	kubernetesClient kubernetes.KubernetesClient
}

func SetupCommand() Command {
	return Command{
		Address:           os.Getenv(VaultAddress),
		AuthToken:         os.Getenv(VaultToken),
		VaultNamespace:    os.Getenv(VaultNamespace),
		EngineName:        os.Getenv(VaultEngine),
		SecretPath:        os.Getenv(VaultSecretPath),
		Base64Kubeconfig:  os.Getenv(Kubeconfig),
		Namespace:         os.Getenv(Namespace),
		ObjectNameToApply: os.Getenv(ObjectNameToApply),
		LoadAsConfigMap:   os.Getenv(ApplyAsConfigmap) == "true",
	}
}

func SetupCommandWithKubernetesClient(args map[string]string, kubernetesClient kubernetes.KubernetesClient) Command {
	_ = os.Setenv(VaultAddress, args[VaultAddress])
	_ = os.Setenv(VaultToken, args[VaultToken])
	_ = os.Setenv(VaultNamespace, args[VaultNamespace])
	_ = os.Setenv(VaultEngine, args[VaultEngine])
	_ = os.Setenv(VaultSecretPath, args[VaultSecretPath])
	_ = os.Setenv(Kubeconfig, args[Kubeconfig])
	_ = os.Setenv(Namespace, args[Namespace])
	_ = os.Setenv(ApplyAsConfigmap, args[ApplyAsConfigmap])
	_ = os.Setenv(ObjectNameToApply, args[ObjectNameToApply])

	command := SetupCommand()
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

	err = kubernetesClient.ApplySecret(context.TODO(), command.ObjectNameToApply, data, log)
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

	err = kubernetesClient.ApplyConfigMap(context.TODO(), command.ObjectNameToApply, data, log)
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
