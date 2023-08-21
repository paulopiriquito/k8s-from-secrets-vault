package tests

import (
	vaultclient "k8s-from-secrets-vault/vault"
	"net"
	"strings"
	"testing"
)

func Test_VaultClient_GivenVaultConfig_CanLoadSecretData(t *testing.T) {
	//Arrange
	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{
		"TEST_KEY": "TEST_VALUE",
	}

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecrets(t, expectedSecretData)
	defer destroyVaultHttpListener(t, vaultHttpListener)

	//Act
	secretData, err := vaultclient.LoadSecretData(clientConfig, log)

	//Assert
	if err != nil {
		t.Error(err)
	}
	if secretData == nil {
		t.Error("Expected secret data to be loaded")
	}
	for key, value := range expectedSecretData {
		if secretData[key] != value {
			t.Errorf("Expected secret data to contain value for %s", key)
		}
	}
}

func Test_VaultClient_GivenEmptySecretTestData_LoadsEmptyMap(t *testing.T) {
	//Arrange
	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{}

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecrets(t, expectedSecretData)
	defer destroyVaultHttpListener(t, vaultHttpListener)

	//Act
	secretData, err := vaultclient.LoadSecretData(clientConfig, log)

	//Assert
	if err != nil {
		t.Error(err)
	}
	if secretData == nil {
		t.Error("Expected secret data to be loaded")
	}
	for key, value := range expectedSecretData {
		if secretData[key] != value {
			t.Errorf("Expected secret data to contain value for %s", key)
		}
	}
}

func Test_VaultClient_GivenIncorrectVaultPath_And_CorrectEngine_LoadsEmptyMap(t *testing.T) {
	//Arrange
	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{}

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecrets(t, expectedSecretData)
	defer destroyVaultHttpListener(t, vaultHttpListener)

	clientConfig.SecretPath = "incorrect/path"

	//Act
	secretData, err := vaultclient.LoadSecretData(clientConfig, log)

	//Assert
	if err != nil {
		t.Error("Expected no error")
	}
	if secretData == nil {
		t.Error("Expected secret data to not be nil")
	}
	if len(secretData) > 0 {
		t.Error("Expected secret data to be empty")
	}
}

func Test_VaultClient_GivenIncorrectEngineName_ReturnsError(t *testing.T) {
	//Arrange
	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{}

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecrets(t, expectedSecretData)
	defer destroyVaultHttpListener(t, vaultHttpListener)

	clientConfig.EngineName = "incorrect-engine-name"

	//Act
	_, err := vaultclient.LoadSecretData(clientConfig, log)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "vault engine does not exist" {
		t.Error("Expected error to be 'vault engine does not exist'")
	}
}

func Test_VaultClient_GivenIncorrectVaultAddress_ReturnsError(t *testing.T) {
	//Arrange
	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{}

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecrets(t, expectedSecretData)
	defer destroyVaultHttpListener(t, vaultHttpListener)

	clientConfig.Address = "http://incorrect-address"

	//Act
	_, err := vaultclient.LoadSecretData(clientConfig, log)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if !strings.Contains(err.Error(), "incorrect-address: no such host") {
		t.Error("Expected error to be 'incorrect-address: no such host'")
	}
}

func Test_VaultClient_GivenIncorrectAuthToken_ReturnsError(t *testing.T) {
	//Arrange
	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{}

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecrets(t, expectedSecretData)
	defer destroyVaultHttpListener(t, vaultHttpListener)

	clientConfig.AuthToken = "incorrect-auth-token"

	//Act
	_, err := vaultclient.LoadSecretData(clientConfig, log)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if !strings.Contains(err.Error(), "permission denied") {
		t.Error("Expected error to be 'permission denied'")
	}
}

func Test_VaultClient_GivenEmptyConfiguration_ReturnsError(t *testing.T) {
	//Arrange
	log := setupLogger(t)

	clientConfig := vaultclient.VaultConfig{}

	//Act
	_, err := vaultclient.LoadSecretData(clientConfig, log)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if !strings.Contains(err.Error(), "address is required") {
		t.Error("Expected error to be 'address is required'")
	}
}

func getClientConfigForNewTestVaultWithSecrets(t *testing.T, secretsToWrite map[string]interface{}) (vaultclient.VaultConfig, net.Listener) {
	t.Helper()

	testVaultConfig := vaultclient.VaultConfig{
		EngineName: "application",
		SecretPath: "dev/config",
		Namespace:  "",
	}

	secretPath := vaultclient.GetSecretPath(testVaultConfig)

	if secretsToWrite == nil {
		secretsToWrite = map[string]interface{}{}
	}

	vaultHttpListener, testVaultRootToken, testVaultAddress := createTestVaultWithSecrets(t, testVaultConfig, secretPath, secretsToWrite)

	clientConfig := vaultclient.VaultConfig{}
	clientConfig.EngineName = testVaultConfig.EngineName
	clientConfig.Namespace = testVaultConfig.Namespace
	clientConfig.SecretPath = testVaultConfig.SecretPath
	clientConfig.Address = testVaultAddress
	clientConfig.AuthToken = testVaultRootToken

	return clientConfig, vaultHttpListener
}
