package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	kubernetes "k8s-from-secrets-vault/kubernetes"
	vault "k8s-from-secrets-vault/vault"
	"os"
)

const (
	VaultAddress         = "VAULT_ADDRESS"
	VaultAuthMethod      = "VAULT_AUTH_METHOD"
	GithubToken          = "GITHUB_TOKEN"
	VaultAppRoleId       = "VAULT_APPROLE_ID"
	VaultAppRoleSecretId = "VAULT_APPROLE_SECRET_ID"
	VaultToken           = "VAULT_TOKEN"
	VaultNamespace       = "VAULT_NAMESPACE"
	VaultEngine          = "VAULT_ENGINE"
	VaultSecretPath      = "VAULT_SECRET_PATH"
	Kubeconfig           = "KUBECONFIG"
	Namespace            = "KUBERNETES_NAMESPACE"
	ApplyAsConfigmap     = "LOAD_AS_CONFIGMAP"
	ObjectNameToApply    = "OBJECT_NAME_TO_APPLY"
)

type Command struct {
	Address         string
	AuthToken       string
	AuthMethod      string
	AppRoleId       string
	AppRoleSecretId string
	GithubToken     string
	VaultNamespace  string
	EngineName      string
	SecretPath      string

	Base64Kubeconfig string
	Namespace        string

	LoadAsConfigMap   bool
	ObjectNameToApply string

	kubernetesClient kubernetes.KubernetesClient
}

func SetupCommand() (*Command, error) {
	log := setupLogger()

	var command = Command{
		Address:           os.Getenv(VaultAddress),
		AuthToken:         os.Getenv(VaultToken),
		AuthMethod:        os.Getenv(VaultAuthMethod),
		AppRoleId:         os.Getenv(VaultAppRoleId),
		AppRoleSecretId:   os.Getenv(VaultAppRoleSecretId),
		GithubToken:       os.Getenv(GithubToken),
		VaultNamespace:    os.Getenv(VaultNamespace),
		EngineName:        os.Getenv(VaultEngine),
		SecretPath:        os.Getenv(VaultSecretPath),
		Base64Kubeconfig:  os.Getenv(Kubeconfig),
		Namespace:         os.Getenv(Namespace),
		ObjectNameToApply: os.Getenv(ObjectNameToApply),
		LoadAsConfigMap:   os.Getenv(ApplyAsConfigmap) == "true",
	}

	if command.AuthMethod == "" {
		command.AuthMethod = "token"
	}

	err := command.Validate()
	if err != nil {
		log.WithError(err).Error("Failed to validate command")
		return nil, err
	}

	return &command, nil
}

func SetupCommandWithKubernetesClient(args map[string]string, kubernetesClient kubernetes.KubernetesClient) (*Command, error) {
	_ = os.Setenv(VaultAddress, args[VaultAddress])
	_ = os.Setenv(VaultToken, args[VaultToken])
	_ = os.Setenv(VaultAuthMethod, args[VaultAuthMethod])
	_ = os.Setenv(GithubToken, args[GithubToken])
	_ = os.Setenv(VaultNamespace, args[VaultNamespace])
	_ = os.Setenv(VaultEngine, args[VaultEngine])
	_ = os.Setenv(VaultSecretPath, args[VaultSecretPath])
	_ = os.Setenv(Kubeconfig, args[Kubeconfig])
	_ = os.Setenv(Namespace, args[Namespace])
	_ = os.Setenv(ApplyAsConfigmap, args[ApplyAsConfigmap])
	_ = os.Setenv(ObjectNameToApply, args[ObjectNameToApply])
	_ = os.Setenv(VaultAppRoleId, args[VaultAppRoleId])
	_ = os.Setenv(VaultAppRoleSecretId, args[VaultAppRoleSecretId])

	command, err := SetupCommand()
	if err != nil {
		return nil, err
	}

	command.kubernetesClient = kubernetesClient
	return command, nil
}

func (command Command) Execute() error {
	if command.LoadAsConfigMap {
		return command.loadAndApplyConfigMap()
	}
	return command.loadAndApplySecrets()
}

func (command Command) vaultParameters() vault.VaultConfig {
	return vault.VaultConfig{
		Address:     command.Address,
		AuthMethod:  command.AuthMethod,
		GithubToken: command.GithubToken,
		AppRoleId:   command.AppRoleId,
		SecretId:    command.AppRoleSecretId,
		AuthToken:   command.AuthToken,
		Namespace:   command.VaultNamespace,
		EngineName:  command.EngineName,
		SecretPath:  command.SecretPath,
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

func (command Command) Validate() error {
	if command.Address == "" {
		return NewError("Vault address is required")
	}
	if command.EngineName == "" {
		return NewError("Vault engine name is required")
	}
	if command.SecretPath == "" {
		return NewError("Vault secret path is required")
	}

	if command.AuthMethod == "approle" && (command.AppRoleId == "" || command.AppRoleSecretId == "") {
		return NewError("Vault RoleId and SecretId are required")
	}
	if command.AuthMethod == "github" && command.GithubToken == "" {
		return NewError("Github token is required")
	}
	if command.AuthMethod == "token" && command.AuthToken == "" {
		return NewError("Vault token is required")
	}

	if command.Base64Kubeconfig == "" {
		return NewError("Kubeconfig is required")
	}
	if command.Namespace == "" {
		return NewError("Kubernetes namespace is required")
	}
	if command.ObjectNameToApply == "" {
		return NewError("Kubernetes object name to apply is required")
	}
	return nil
}

func NewError(s string) error {
	return fmt.Errorf(s)
}

func setupLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{}
	return log
}
