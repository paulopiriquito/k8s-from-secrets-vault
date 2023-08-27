package vault_client

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

type VaultConfig struct {
	Address     string
	AuthToken   string
	Namespace   string
	EngineName  string
	SecretPath  string
	AuthMethod  string
	GithubToken string
	AppRoleId   string
	SecretId    string
}

func CheckVaultConfigRequiredFields(config VaultConfig) error {
	if config.Address == "" {
		return fmt.Errorf("address is required")
	}
	if config.AuthMethod == "token" && config.AuthToken == "" {
		return fmt.Errorf("authToken is required")
	}
	if config.AuthMethod == "github" && config.GithubToken == "" {
		return fmt.Errorf("githubToken is required")
	}
	if config.EngineName == "" {
		return fmt.Errorf("engineName is required")
	}
	return nil
}

func LoadSecretData(config VaultConfig, log *logrus.Logger) (map[string]string, error) {
	log.WithFields(logrus.Fields{
		"address":    config.Address,
		"namespace":  config.Namespace,
		"secretPath": GetSecretPath(config),
	}).Info("Loading secret data")

	err := CheckVaultConfigRequiredFields(config)
	if err != nil {
		return nil, err
	}

	client, err := newAuthenticatedVaultApiClient(config, log)
	if err != nil {
		return nil, err
	}

	secret, err := client.Logical().Read(GetSecretPath(config))
	if err != nil {
		log.Error(err, "Failed to read vault engine path %s", GetSecretPath(config))
		return nil, err
	}

	if secret == nil {
		if !vaultEngineExists(config, log, client) {
			log.Error(nil, "Vault engine does not exist")
			return nil, fmt.Errorf("vault engine does not exist")
		}
		log.Warning("Secret engine path is empty")
		return map[string]string{}, nil
	}

	if secret.Data == nil {
		log.Error(nil, "Failed to parse Vault secret data")
		return nil, fmt.Errorf("failed to parse Vault secret data")
	}

	secretData := make(map[string]string)

	// If secret.Data has key named "data" then it is a KVv2 secret
	if _, ok := secret.Data["data"]; ok {
		for k, v := range secret.Data["data"].(map[string]interface{}) {
			var value = ""
			if v != nil {
				value = fmt.Sprintf("%v", v)
			}
			secretData[k] = value
		}
	} else {
		for k, v := range secret.Data {
			var value = ""
			if v != nil {
				value = fmt.Sprintf("%v", v)
			}
			secretData[k] = value
		}
	}

	log.WithFields(logrus.Fields{
		"address":    config.Address,
		"namespace":  config.Namespace,
		"secretPath": GetSecretPath(config),
	}).Info("Loaded secret data")

	return secretData, nil
}

func GetSecretPath(config VaultConfig) string {
	return fmt.Sprintf("%s/data/%s", config.EngineName, config.SecretPath)
}

func vaultEngineExists(config VaultConfig, log *logrus.Logger, client *api.Client) bool {
	mounts, err := client.Sys().ListMounts()

	if err != nil {
		log.Error(err, "Failed to fetch Vault secret engine mounts")
		return false
	}
	if _, ok := mounts[config.EngineName+"/"]; !ok {
		return false
	}

	return true
}

func newAuthenticatedVaultApiClient(config VaultConfig, log *logrus.Logger) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: config.Address,
	})

	if err != nil {
		log.WithError(err).Error("Failed to create vault client")
		return nil, err
	}

	client.SetNamespace(config.Namespace)

	if config.AuthMethod == "approle" {
		client, err = authWithAppRole(config.AppRoleId, config.SecretId, client)
		if err != nil {
			log.WithError(err).Error("Failed to authenticate with AppRole")
			return nil, err
		}
	}
	if config.AuthMethod == "github" {
		client, err = authWithGithub(config.GithubToken, client)
		if err != nil {
			log.WithError(err).Error("Failed to authenticate with Github")
			return nil, err
		}
	}
	if config.AuthMethod == "token" {
		client, err = authWithToken(config.AuthToken, client)
		if err != nil {
			log.WithError(err).Error("Failed to authenticate with token")
			return nil, err
		}
	}

	return client, nil
}

func authWithAppRole(roleId string, secretId string, client *api.Client) (*api.Client, error) {
	secret, err := client.Logical().Write(fmt.Sprintf("%s/auth/approle/login", client.Namespace()), map[string]interface{}{
		"role_id":   roleId,
		"secret_id": secretId,
	})

	if err != nil {
		return client, err
	}
	if secret == nil {
		return client, nil
	}
	if secret.Auth == nil {
		return client, fmt.Errorf("secret.Auth is nil")
	}
	if secret.Auth.ClientToken == "" {
		return client, fmt.Errorf("secret.Auth.ClientToken is empty")
	}

	client.SetToken(secret.Auth.ClientToken)

	return client, nil
}

func authWithToken(token string, client *api.Client) (*api.Client, error) {
	client.SetToken(token)
	return client, nil
}

func authWithGithub(githubToken string, client *api.Client) (*api.Client, error) {
	secret, err := client.Logical().Write(fmt.Sprintf("%s/auth/github/login", client.Namespace()), map[string]interface{}{
		"token": githubToken,
	})

	if err != nil {
		return client, err
	}
	if secret == nil {
		return client, nil
	}
	if secret.Auth == nil {
		return client, fmt.Errorf("secret.Auth is nil")
	}
	if secret.Auth.ClientToken == "" {
		return client, fmt.Errorf("secret.Auth.ClientToken is empty")
	}

	client.SetToken(secret.Auth.ClientToken)

	return client, nil
}
