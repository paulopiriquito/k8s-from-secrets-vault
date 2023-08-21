package vault_client

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

type VaultConfig struct {
	Address    string
	AuthToken  string
	Namespace  string
	EngineName string
	SecretPath string
}

func CheckVaultConfigRequiredFields(config VaultConfig) error {
	if config.Address == "" {
		return fmt.Errorf("address is required")
	}
	if config.AuthToken == "" {
		return fmt.Errorf("authToken is required")
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
	for k, v := range secret.Data {
		secretData[k] = v.(string)
	}

	log.WithFields(logrus.Fields{
		"address":    config.Address,
		"namespace":  config.Namespace,
		"secretPath": GetSecretPath(config),
	}).Info("Loaded secret data")

	return secretData, nil
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

func GetSecretPath(config VaultConfig) string {
	return fmt.Sprintf("%s/data/%s", config.EngineName, config.SecretPath)
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
	client.SetToken(config.AuthToken)

	return client, nil
}
