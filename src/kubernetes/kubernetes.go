package kubernetes_client

import (
	"context"
	"encoding/base64"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type KubernetesClient interface {
	ApplySecret(context context.Context, secretName string, secretData map[string]string) error
	ApplyConfigMap(context context.Context, configName string, configData map[string]string) error
}
type kubernetesClient struct {
	config KubernetesConfig
	client *kubernetes.Clientset
}

type KubernetesConfig struct {
	restConfig *rest.Config
	namespace  string
}

type KubernetesParameters struct {
	Base64Kubeconfig string
	Namespace        string
}

func NewKubernetesClient(config KubernetesConfig) (KubernetesClient, error) {
	client, err := kubernetes.NewForConfig(config.restConfig)
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
		return nil, err
	}

	return kubernetesClient{config, client}, nil
}

func NewKubernetesConfig(kubernetesParameters KubernetesParameters) (KubernetesConfig, error) {
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubernetesParameters.Base64Kubeconfig)
	if err != nil {
		log.Fatalf("Error decoding base64 kubeconfig: %v", err)
		return KubernetesConfig{}, err
	}

	// Create a Kubernetes client configuration from the decoded kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)

	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
		return KubernetesConfig{}, err
	}

	return KubernetesConfig{restConfig: config, namespace: kubernetesParameters.Namespace}, nil
}

func (c kubernetesClient) ApplySecret(context context.Context, secretName string, secretData map[string]string) error {
	secret := coreV1.Secret(secretName, c.config.namespace)

	secret = secret.WithType("Opaque")
	secret = secret.WithStringData(secretData)
	secret = secret.WithImmutable(true)

	_, err := c.client.CoreV1().Secrets(c.config.namespace).Apply(context, secret, metav1.ApplyOptions{})
	if err != nil {
		log.Fatalf("Error applying Secret: %v", err)
		return err
	}

	return nil
}

func (c kubernetesClient) ApplyConfigMap(context context.Context, configName string, configData map[string]string) error {
	configMap := coreV1.ConfigMap(configName, c.config.namespace)

	configMap = configMap.WithData(configData)
	configMap = configMap.WithImmutable(true)

	_, err := c.client.CoreV1().ConfigMaps(c.config.namespace).Apply(context, configMap, metav1.ApplyOptions{})
	if err != nil {
		log.Fatalf("Error applying ConfigMap: %v", err)
		return err
	}

	return nil
}
