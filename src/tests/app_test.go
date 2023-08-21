package tests

import (
	"context"
	"k8s-from-secrets-vault/app"
	kubernetes "k8s-from-secrets-vault/kubernetes"
	vaultclient "k8s-from-secrets-vault/vault"
	"strings"
	"testing"
)

func Test_Command_GivenRequiredArgs_CanLoadAndApplySecrets(t *testing.T) {
	command := app.Command{
		Address:          "http://",
		AuthToken:        "",
		VaultNamespace:   "",
		EngineName:       "",
		SecretPath:       "",
		Base64Kubeconfig: "",
		Namespace:        "",
	}

	err := command.LoadAndApplySecrets()

	if err != nil {
		t.Error("Expected no error, got ", err)
	}
}

func Test_Command_GivenIncorrectArgs_CanHandleError(t *testing.T) {
	command := app.Command{
		Address:          "http://",
		AuthToken:        "",
		VaultNamespace:   "",
		EngineName:       "",
		SecretPath:       "",
		Base64Kubeconfig: "",
		Namespace:        "",
	}

	err := command.LoadAndApplyConfigMap()

	if err != nil {
		t.Error("Expected no error, got ", err)
	}
}

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
	if !strings.Contains(err.Error(), "incorrect-address") {
		t.Error("Expected error to be 'incorrect-address'")
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

func Test_KubernetesClient_GivenKubernetesConfigAndSecretData_CanApplySecret(t *testing.T) {
	kubernetesConfig, err := kubernetes.NewKubernetesConfig(kubernetes.KubernetesParameters{Namespace: "default"})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	kubernetesClient, err := kubernetes.NewKubernetesClient(kubernetesConfig)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = kubernetesClient.ApplySecret(context.TODO(), "app-config", map[string]string{
		"key": "value",
	})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
}

func Test_KubernetesClient_GivenKubernetesConfigAndSecretData_CanApplyConfigMap(t *testing.T) {
	kubernetesConfig, err := kubernetes.NewKubernetesConfig(kubernetes.KubernetesParameters{Namespace: "default"})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	kubernetesClient, err := kubernetes.NewKubernetesClient(kubernetesConfig)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = kubernetesClient.ApplyConfigMap(context.TODO(), "app-config", map[string]string{
		"key": "value",
	})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
}

func Test_KubernetesClient_GivenEmptySecretData_CanApplyEmptySecret(t *testing.T) {
	kubernetesConfig, err := kubernetes.NewKubernetesConfig(kubernetes.KubernetesParameters{Namespace: "default"})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	kubernetesClient, err := kubernetes.NewKubernetesClient(kubernetesConfig)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = kubernetesClient.ApplySecret(context.TODO(), "app-secret", map[string]string{})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
}

func Test_KubernetesClient_GivenEmptySecretData_CanApplyEmptyConfigMap(t *testing.T) {
	kubernetesConfig, err := kubernetes.NewKubernetesConfig(kubernetes.KubernetesParameters{Namespace: "default"})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	kubernetesClient, err := kubernetes.NewKubernetesClient(kubernetesConfig)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = kubernetesClient.ApplyConfigMap(context.TODO(), "app-config", map[string]string{})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
}

func Test_KubernetesClient_GivenIncorrectKubernetesConfig_CanHandleError(t *testing.T) {
	t.Skip("Not implemented")
}
