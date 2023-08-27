package tests

import (
	"context"
	"k8s-from-secrets-vault/app"
	kubernetes "k8s-from-secrets-vault/kubernetes"
	vaultclient "k8s-from-secrets-vault/vault"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
)

func Test_GivenMissingRequiredObjectNameToApply_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:     "http://",
		app.VaultToken:       "test-token",
		app.VaultEngine:      "test-engine",
		app.VaultSecretPath:  "test-path",
		app.Namespace:        "test-namespace",
		app.ApplyAsConfigmap: "false",
		app.Kubeconfig:       "test-kubeconfig",
		app.VaultAuthMethod:  "jwt",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Kubernetes object name to apply is required" {
		t.Error("Expected error to be 'kubernetes object name to apply is required'")
	}
}

func Test_GivenMissingRequiredNamespace_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:      "http://",
		app.VaultToken:        "test-token",
		app.VaultEngine:       "test-engine",
		app.VaultSecretPath:   "test-path",
		app.ApplyAsConfigmap:  "false",
		app.Kubeconfig:        "test-kubeconfig",
		app.ObjectNameToApply: "test-secret",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Kubernetes namespace is required" {
		t.Error("Expected error to be 'kubernetes namespace is required'")
	}
}

func Test_WhenAuthMethodIsGithub_GivenMissingRequiredGithubToken_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:      "http://",
		app.VaultAuthMethod:   "github",
		app.VaultEngine:       "test-engine",
		app.VaultSecretPath:   "test-path",
		app.Namespace:         "test-namespace",
		app.ApplyAsConfigmap:  "false",
		app.Kubeconfig:        "test-kubeconfig",
		app.ObjectNameToApply: "test-secret",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Github token is required" {
		t.Error("Expected error to be 'github token is required'")
	}
}

func Test_WhenAuthMethodAppRole_GivenMissingRequiredAppRoleAndSecretId_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:         "http://",
		app.VaultEngine:          "test-engine",
		app.VaultSecretPath:      "test-path",
		app.Namespace:            "test-namespace",
		app.ApplyAsConfigmap:     "false",
		app.Kubeconfig:           "test-kubeconfig",
		app.ObjectNameToApply:    "test-secret",
		app.VaultAuthMethod:      "approle",
		app.VaultAppRoleId:       "",
		app.VaultAppRoleSecretId: "",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Vault RoleId and SecretId are required" {
		t.Error("Expected error to be 'vault RoleId and SecretId are required'")
	}
}

func Test_GivenMissingRequiredVaultSecretPath_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:      "http://",
		app.VaultAuthMethod:   "jwt",
		app.VaultToken:        "test-token",
		app.VaultEngine:       "test-engine",
		app.Namespace:         "test-namespace",
		app.ApplyAsConfigmap:  "false",
		app.Kubeconfig:        "test-kubeconfig",
		app.ObjectNameToApply: "test-secret",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Vault secret path is required" {
		t.Error("Expected error to be 'vault secret path is required'")
	}
}

func Test_GivenMissingRequiredVaultToken_WhenAuthMethodJwt_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:      "http://",
		app.VaultEngine:       "test-engine",
		app.VaultSecretPath:   "test-path",
		app.Namespace:         "test-namespace",
		app.ApplyAsConfigmap:  "false",
		app.Kubeconfig:        "test-kubeconfig",
		app.VaultAuthMethod:   "jwt",
		app.ObjectNameToApply: "test-secret",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Vault token is required" {
		t.Error("Expected error to be 'vault token is required'")
	}
}

func Test_GivenMissingRequiredVaultAddress_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAuthMethod:   "jwt",
		app.VaultToken:        "test-token",
		app.VaultEngine:       "test-engine",
		app.VaultSecretPath:   "test-path",
		app.Namespace:         "test-namespace",
		app.ApplyAsConfigmap:  "false",
		app.Kubeconfig:        "test-kubeconfig",
		app.ObjectNameToApply: "test-secret",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Vault address is required" {
		t.Error("Expected error to be 'vault address is required'")
	}
}

func Test_GivenMissingRequiredKubeconfig_ReturnsError(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:      "http://",
		app.VaultAuthMethod:   "jwt",
		app.VaultToken:        "test-token",
		app.VaultEngine:       "test-engine",
		app.VaultSecretPath:   "test-path",
		app.Namespace:         "test-namespace",
		app.ApplyAsConfigmap:  "false",
		app.ObjectNameToApply: "test-secret",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Kubeconfig is required" {
		t.Error("Expected error to be 'kubeconfig is required'")
	}
}

func Test_GivenMissingRequiredVaultEngine_ReturnsError(t *testing.T) {
	//Arrange missing vault engine
	commandArgs := map[string]string{
		app.Kubeconfig:        "test-kubeconfig",
		app.Namespace:         "test-namespace",
		app.ApplyAsConfigmap:  "false",
		app.ObjectNameToApply: "test-secret",
		app.VaultAddress:      "http://",
		app.VaultSecretPath:   "test-path",
		app.VaultAuthMethod:   "jwt",
		app.VaultToken:        "test-token",
	}

	//Act
	_, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Vault engine name is required" {
		t.Error("Expected error to be 'vault engine is required'")
	}
}

func Test_WhenAuthMethodIsNotSet_DefaultsToJwt(t *testing.T) {
	//Arrange
	commandArgs := map[string]string{
		app.VaultAddress:      "http://",
		app.VaultToken:        "test-token",
		app.VaultEngine:       "test-engine",
		app.VaultSecretPath:   "test-path",
		app.Namespace:         "test-namespace",
		app.ApplyAsConfigmap:  "false",
		app.Kubeconfig:        "test-kubeconfig",
		app.ObjectNameToApply: "test-secret",
	}

	//Act
	command, err := app.SetupCommandWithKubernetesClient(commandArgs, nil)

	//Assert
	if err != nil {
		t.Error("Expected no error")
	}
	if command.AuthToken != "test-token" {
		t.Error("Expected command.AuthToken to be 'test-token'")
	}
}

func Test_Command_GivenRequiredArgs_CanLoadAndApplySecrets(t *testing.T) {
	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{
		"TEST_KEY": "TEST_VALUE",
	}

	vaultClientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecretsAndJwtAuth(t, expectedSecretData)
	defer destroyVaultHttpListener(t, vaultHttpListener)

	parameters := getFakeKubernetesParameters(t)

	//Transform command into a string array
	commandArgs := map[string]string{
		app.VaultAddress:      vaultClientConfig.Address,
		app.VaultToken:        vaultClientConfig.AuthToken,
		app.VaultEngine:       vaultClientConfig.EngineName,
		app.VaultSecretPath:   vaultClientConfig.SecretPath,
		app.Kubeconfig:        parameters.Base64Kubeconfig,
		app.Namespace:         parameters.Namespace,
		app.VaultAuthMethod:   "jwt",
		app.ApplyAsConfigmap:  "false",
		app.ObjectNameToApply: "test-secret",
	}

	config, err := kubernetes.CreateConfig(parameters, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	client, fakeClient, err := getFakeKubernetesClient(config, t)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	command, err := app.SetupCommandWithKubernetesClient(commandArgs, client)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	err = command.Execute()
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	secret, err := fakeClient.CoreV1().Secrets("test-namespace").Get(context.TODO(), "test-secret", metav1.GetOptions{})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if secret.StringData["TEST_KEY"] != "TEST_VALUE" {
		t.Error("Expected secret to contain TEST_KEY with value TEST_VALUE")
	}
}

func Test_Command_GivenIncorrectArgs_ReturnsError(t *testing.T) {
	command := app.Command{
		Address:          "http://",
		AuthToken:        "",
		VaultNamespace:   "",
		EngineName:       "",
		SecretPath:       "",
		Base64Kubeconfig: "",
		Namespace:        "",
		LoadAsConfigMap:  false,
	}

	err := command.Execute()

	if err == nil {
		t.Error("Expected error")
	}
}

func Test_VaultClient_GivenVaultConfig_CanLoadSecretData(t *testing.T) {
	//Arrange

	log := setupLogger(t)

	expectedSecretData := map[string]interface{}{
		"TEST_KEY": "TEST_VALUE",
	}

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecretsAndJwtAuth(t, expectedSecretData)
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

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecretsAndJwtAuth(t, expectedSecretData)
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

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecretsAndJwtAuth(t, expectedSecretData)
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

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecretsAndJwtAuth(t, expectedSecretData)
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

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecretsAndJwtAuth(t, expectedSecretData)
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

	clientConfig, vaultHttpListener := getClientConfigForNewTestVaultWithSecretsAndJwtAuth(t, expectedSecretData)
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
	log := setupLogger(t)

	parameters := kubernetes.KubernetesParameters{
		Base64Kubeconfig: getFakeKubeconfigAsBase64(t),
		Namespace:        "test",
	}

	config, err := kubernetes.CreateConfig(parameters, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	client, fakeClient, err := getFakeKubernetesClient(config, t)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = client.ApplySecret(context.Background(), "test-secret", map[string]string{"TEST_KEY": "TEST_VALUE"}, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	secret, err := fakeClient.CoreV1().Secrets("test").Get(context.Background(), "test-secret", metav1.GetOptions{})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if secret.StringData["TEST_KEY"] != "TEST_VALUE" {
		t.Error("Expected secret to contain TEST_KEY with value TEST_VALUE")
	}
}

func Test_KubernetesClient_GivenKubernetesConfigAndSecretData_CanApplyConfigMap(t *testing.T) {
	log := setupLogger(t)

	parameters := kubernetes.KubernetesParameters{
		Base64Kubeconfig: getFakeKubeconfigAsBase64(t),
		Namespace:        "test",
	}

	config, err := kubernetes.CreateConfig(parameters, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	client, fakeClient, err := getFakeKubernetesClient(config, t)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = client.ApplyConfigMap(context.Background(), "test-config", map[string]string{"TEST_KEY": "TEST_VALUE"}, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	configMap, err := fakeClient.CoreV1().ConfigMaps("test").Get(context.Background(), "test-config", metav1.GetOptions{})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if configMap.Data["TEST_KEY"] != "TEST_VALUE" {
		t.Error("Expected config map to contain TEST_KEY with value TEST_VALUE")
	}
}

func Test_KubernetesClient_GivenEmptySecretData_CanApplyEmptySecret(t *testing.T) {
	log := setupLogger(t)

	parameters := kubernetes.KubernetesParameters{
		Base64Kubeconfig: getFakeKubeconfigAsBase64(t),
		Namespace:        "test",
	}

	config, err := kubernetes.CreateConfig(parameters, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	client, fakeClient, err := getFakeKubernetesClient(config, t)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = client.ApplySecret(context.Background(), "test-secret", map[string]string{}, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	secret, err := fakeClient.CoreV1().Secrets("test").Get(context.Background(), "test-secret", metav1.GetOptions{})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if len(secret.StringData) > 0 {
		t.Error("Expected secret to be empty")
	}
}

func Test_KubernetesClient_GivenEmptySecretData_CanApplyEmptyConfigMap(t *testing.T) {
	log := setupLogger(t)

	parameters := kubernetes.KubernetesParameters{
		Base64Kubeconfig: getFakeKubeconfigAsBase64(t),
		Namespace:        "test",
	}

	config, err := kubernetes.CreateConfig(parameters, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	client, fakeClient, err := getFakeKubernetesClient(config, t)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	err = client.ApplyConfigMap(context.Background(), "test-config", map[string]string{}, log)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	configMap, err := fakeClient.CoreV1().ConfigMaps("test").Get(context.Background(), "test-config", metav1.GetOptions{})
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if len(configMap.Data) > 0 {
		t.Error("Expected config map to be empty")
	}
}

func Test_KubernetesClient_GivenIncorrectKubernetesParameters_ReturnsErrorWhenCreatingConfig(t *testing.T) {
	log := setupLogger(t)

	parameters := kubernetes.KubernetesParameters{}

	_, err := kubernetes.CreateConfig(parameters, log)

	if err == nil {
		t.Error("Expected error")
	}
}

func Test_KubernetesClient_GivenValidParameters_CanCreateConfig(t *testing.T) {
	log := setupLogger(t)

	parameters := kubernetes.KubernetesParameters{
		Base64Kubeconfig: getFakeKubeconfigAsBase64(t),
		Namespace:        "test",
	}

	config, err := kubernetes.CreateConfig(parameters, log)

	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if config.GetNamespace() != "test" {
		t.Error("Expected namespace to be 'test'")
	}
	if config.GetServer() != "https://example.com" {
		t.Error("Expected server to be 'https://example.com'")
	}
}
