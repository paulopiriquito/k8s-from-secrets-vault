package kubernetes_client

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient interface {
	ApplySecret(context context.Context, secretName string, secretData map[string]string, log *logrus.Logger) error
	ApplyConfigMap(context context.Context, configName string, configData map[string]string, log *logrus.Logger) error
}
type kubernetesClient struct {
	config     KubernetesConfig
	client     kubernetes.Interface
	commitMode bool
}

type KubernetesConfig struct {
	restConfig *rest.Config
	namespace  string
}

type KubernetesParameters struct {
	Base64Kubeconfig string
	Namespace        string
}

const CREATE = true
const APPLY = false

func InjectKubernetesClient(client kubernetes.Interface, config KubernetesConfig) KubernetesClient {
	return kubernetesClient{config, client, CREATE}
}

func CreateClient(config KubernetesConfig, log *logrus.Logger) (KubernetesClient, error) {
	client, err := kubernetes.NewForConfig(config.restConfig)
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
		return nil, err
	}

	return kubernetesClient{config, client, APPLY}, nil
}

func CreateConfig(kubernetesParameters KubernetesParameters, log *logrus.Logger) (KubernetesConfig, error) {
	log.WithFields(logrus.Fields{
		"namespace": kubernetesParameters.Namespace,
	}).Info("Loading kubeconfig data")

	if kubernetesParameters.Base64Kubeconfig == "" {
		return KubernetesConfig{}, fmt.Errorf("provided base64 kubeconfig is empty")
	}
	if kubernetesParameters.Namespace == "" {
		return KubernetesConfig{}, fmt.Errorf("provided namespace is empty")
	}

	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubernetesParameters.Base64Kubeconfig)
	if err != nil {
		log.Errorf("Error decoding base64 kubeconfig: %v", err)
		return KubernetesConfig{}, err
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)

	if err != nil {
		log.Errorf("Error building kubeconfig: %v", err)
		return KubernetesConfig{}, err
	}

	return KubernetesConfig{restConfig: config, namespace: kubernetesParameters.Namespace}, nil
}

func (conf KubernetesConfig) GetNamespace() string {
	return conf.namespace
}

func (conf KubernetesConfig) GetServer() string {
	return conf.restConfig.Host
}

func (c kubernetesClient) ApplySecret(context context.Context, secretName string, secretData map[string]string, log *logrus.Logger) error {
	var createdSecret, err = &corev1.Secret{}, error(nil)

	if c.commitMode == CREATE {
		createdSecret, err = c.createSecret(context, secretName, secretData)
	} else if c.commitMode == APPLY {
		createdSecret, err = c.applySecret(context, secretName, secretData)
	}

	if err != nil {
		log.Errorf("Error applying Secret: %v", err)
		return err
	}
	if createdSecret == nil {
		log.Errorf("Error applying Secret: %v", err)
		return err
	}
	return nil
}

func (c kubernetesClient) ApplyConfigMap(context context.Context, configName string, configData map[string]string, log *logrus.Logger) error {
	var createdConfigMap, err = &corev1.ConfigMap{}, error(nil)

	if c.commitMode == CREATE {
		createdConfigMap, err = c.createConfigMap(context, configName, configData)
	} else if c.commitMode == APPLY {
		createdConfigMap, err = c.applyConfigMap(context, configName, configData)
	}

	if err != nil {
		log.Errorf("Error applying Config-Map: %v", err)
		return err
	}
	if createdConfigMap == nil {
		log.Errorf("Error applying Config-Map: %v", err)
		return err
	}
	return nil
}

func (c kubernetesClient) createSecret(context context.Context, secretName string, secretData map[string]string) (*corev1.Secret, error) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: c.config.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/update-by": "k8s-from-secrets-vault",
			},
		},
		StringData: secretData,
		Type:       "Opaque",
	}

	createdSecret, err := c.client.CoreV1().Secrets(c.config.namespace).Create(context, &secret, metav1.CreateOptions{})
	return createdSecret, err
}

func (c kubernetesClient) applySecret(context context.Context, secretName string, secretData map[string]string) (*corev1.Secret, error) {
	secret := applyv1.Secret(secretName, c.config.namespace)
	secret = secret.WithType("Opaque")
	secret = secret.WithStringData(secretData)
	secret = secret.WithAnnotations(map[string]string{
		"app.kubernetes.io/update-by": "k8s-from-secrets-vault",
	})

	appliedSecret, err := c.client.CoreV1().Secrets(c.config.namespace).Apply(context, secret, metav1.ApplyOptions{})
	return appliedSecret, err
}

func (c kubernetesClient) createConfigMap(context context.Context, configName string, configData map[string]string) (*corev1.ConfigMap, error) {
	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: c.config.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/update-by": "k8s-from-secrets-vault",
			},
		},
		Data: configData,
	}

	createdConfigMap, err := c.client.CoreV1().ConfigMaps(c.config.namespace).Create(context, &configmap, metav1.CreateOptions{})
	return createdConfigMap, err
}

func (c kubernetesClient) applyConfigMap(context context.Context, configName string, configData map[string]string) (*corev1.ConfigMap, error) {
	configmap := applyv1.ConfigMap(configName, c.config.namespace)
	configmap = configmap.WithData(configData)
	configmap = configmap.WithAnnotations(map[string]string{
		"app.kubernetes.io/update-by": "k8s-from-secrets-vault",
	})

	appliedConfigMap, err := c.client.CoreV1().ConfigMaps(c.config.namespace).Apply(context, configmap, metav1.ApplyOptions{})
	return appliedConfigMap, err
}
