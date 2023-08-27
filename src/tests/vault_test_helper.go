package tests

import (
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	vaultclient "k8s-from-secrets-vault/vault"
	"net"
	"strings"
	"testing"
)

func getClientConfigForNewTestVaultWithSecretsAndTokenAuth(t *testing.T, secretsToWrite map[string]interface{}) (vaultclient.VaultConfig, net.Listener) {
	t.Helper()
	testVaultConfig := getTestVaultConfigWithAuthMethod("token")
	vaultHttpListener, clientConfig := getClientConfigWithSecrets(t, secretsToWrite, testVaultConfig)
	return clientConfig, vaultHttpListener
}

func getClientConfigForNewTestVaultWithSecretsAndGithubAuth(t *testing.T, secretsToWrite map[string]interface{}) (vaultclient.VaultConfig, net.Listener) {
	t.Helper()
	testVaultConfig := getTestVaultConfigWithAuthMethod("github")
	vaultHttpListener, clientConfig := getClientConfigWithSecrets(t, secretsToWrite, testVaultConfig)
	return clientConfig, vaultHttpListener
}

func getTestVaultConfigWithAuthMethod(authMethod string) vaultclient.VaultConfig {
	testVaultConfig := vaultclient.VaultConfig{
		EngineName: "application",
		SecretPath: "dev/config",
		Namespace:  "",
		AuthMethod: authMethod,
	}
	return testVaultConfig
}

func getClientConfigWithSecrets(t *testing.T, secretsToWrite map[string]interface{}, testVaultConfig vaultclient.VaultConfig) (net.Listener, vaultclient.VaultConfig) {
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
	clientConfig.AuthMethod = testVaultConfig.AuthMethod
	clientConfig.AuthToken = testVaultRootToken
	return vaultHttpListener, clientConfig
}

func createTestVaultWithSecrets(t *testing.T, vaultConfig vaultclient.VaultConfig, secretPath string, testSecrets map[string]interface{}) (net.Listener, string, string) {
	t.Helper()
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{NumCores: 1})
	rootToken := cluster.RootToken
	vaultCore := cluster.Cores[0].Core

	httpServerListener, serverAddress := http.TestServer(t, vaultCore)

	conf := api.DefaultConfig()

	conf.Address = serverAddress

	client := createVaultClient(t, vaultConfig.Namespace, conf, rootToken)

	client = client.WithNamespace(vaultConfig.Namespace)

	createSecretEngineIfMissing(t, client, vaultConfig.EngineName)

	setupTestSecrets(t, client, secretPath, testSecrets)

	return httpServerListener, rootToken, serverAddress
}

func setupTestSecrets(t *testing.T, client *api.Client, secretPath string, testSecrets map[string]interface{}) {
	t.Helper()

	if len(testSecrets) == 0 {
		return
	}

	_, err := client.Logical().Write(secretPath, testSecrets)
	if err != nil {
		t.Fatal(err)
	}
	secret, err := client.Logical().Read(secretPath)
	if err != nil || secret == nil {
		t.Fatal(err)
	}
}

func createSecretEngineIfMissing(t *testing.T, client *api.Client, engineName string) {
	t.Helper()

	err := client.Sys().Mount(engineName, &api.MountInput{Type: "kv"})

	if err != nil && !strings.Contains(err.Error(), "existing mount") {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListMounts()

	if err != nil {
		t.Fatal(err)
	}
	if _, ok := mounts[engineName+"/"]; !ok {
		t.Fatalf("Mount path %s does not exist", engineName)
	}
}

func createVaultClient(t *testing.T, namespace string, conf *api.Config, rootToken string) *api.Client {
	t.Helper()

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)
	client.SetNamespace(namespace)
	return client
}

func destroyVaultHttpListener(t *testing.T, vaultHttpListener net.Listener) {
	err := vaultHttpListener.Close()
	if err != nil {
		t.Error(err)
	}
}
