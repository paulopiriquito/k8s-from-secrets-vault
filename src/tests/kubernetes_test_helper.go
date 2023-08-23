package tests

import (
	"context"
	"encoding/base64"
	kubernetesclient "k8s-from-secrets-vault/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func getFakeKubernetesClient(config kubernetesclient.KubernetesConfig, t *testing.T) (kubernetesclient.KubernetesClient, *fake.Clientset, error) {
	t.Helper()
	fakeClient := fake.NewSimpleClientset()

	namespace := config.GetNamespace()

	namespaceToApply := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, err := fakeClient.CoreV1().Namespaces().Create(context.TODO(), namespaceToApply, metav1.CreateOptions{})
	if err != nil {
		t.Logf("Error creating test namespace %s: %v", namespace, err)
		return nil, nil, err
	}

	return kubernetesclient.InjectKubernetesClient(fakeClient, config), fakeClient, nil
}

func getFakeKubernetesParameters(t *testing.T) kubernetesclient.KubernetesParameters {
	t.Helper()
	return kubernetesclient.KubernetesParameters{
		Base64Kubeconfig: getFakeKubeconfigAsBase64(t),
		Namespace:        "test-namespace",
	}
}

func getFakeKubeconfigAsBase64(t *testing.T) string {
	t.Helper()
	return base64.StdEncoding.EncodeToString([]byte(fakeKubeconfig))
}

const fakeKubeconfig = `
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://example.com
  name: example-cluster
contexts:
- context:
    cluster: example-cluster
    user: example-user
  name: example-context
current-context: example-context
preferences: {}
users:
- name: example-user
  user:
    token: abcdef1234567890
`
