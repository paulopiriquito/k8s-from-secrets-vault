import os
import sys
import json
import requests
import base64
from kubernetes import client, config


def load_kubeconfig(kubeconfig: str):
    # kubeconfig_decoded = base64.b64decode(kubeconfig)
    # with open('./.kube/config', 'w') as f:
    #     f.write(str(kubeconfig_decoded))

    config.load_kube_config(f"./kube/config")
    k8s_apps_client = client.CoreV1Api()
    return k8s_apps_client


def apply_k8s_secret(name: str, secret_data: dict, namespace: str, kubeconfig: str):
    api = load_kubeconfig(kubeconfig)
    create_secret_k8s_request(name, secret_data, namespace, api)


def apply_k8s_configmap(name: str, configmap_data: dict, namespace: str, kubeconfig: str):
    api = load_kubeconfig(kubeconfig)
    create_configmap_k8s_request(name, configmap_data, namespace, api)


def create_secret_k8s_request(name: str, s_data: dict, namespace: str, client_api):
    secret = client.V1Secret(
        api_version="v1",
        kind="Secret",
        metadata=client.V1ObjectMeta(name=name, namespace=namespace),
        data=s_data
    )

    api = client_api.create_namespaced_secret(namespace=namespace, body=secret)
    return api


def create_configmap_k8s_request(name: str, c_data: dict, namespace: str, client_api):
    config_map = client.V1ConfigMap(
        api_version="v1",
        kind="ConfigMap",
        metadata=client.V1ObjectMeta(name=name, namespace=namespace),
        data=c_data
    )

    api = client_api.create_namespaced_config_map(namespace=namespace, body=config_map)
    return api


def remove_special(inputs: dict):
    sk = ['__filename__', '__type__', '__path__']
    new_data = {}
    for k, v in inputs.items():
        if k not in sk:
            new_data[k] = v
    return new_data


def base64encode_collection(collection: dict):
    for k, v in collection:
        data[k] = base64.urlsafe_b64encode(v)
    return collection


# main entrypoint

i_kubeconfig = os.getenv('INPUT_KUBECONFIG')
i_namespace = os.getenv('INPUT_NAMESPACE')
i_resource_name = os.getenv('INPUT_RESOURCE-NAME')
i_resource_type = os.getenv('INPUT_RESOURCE-TYPE')

i_vault = os.getenv('INPUT_VAULT-URL')
i_engine = os.getenv('INPUT_ENGINE-NAME')
i_secret = os.getenv('INPUT_SECRET-NAME')
i_token = os.getenv('INPUT_VAULT-AUTH-TOKEN')

try:
    debug = sys.argv[5]
except:
    debug = 'no'

headers = {
    'X-Vault-Token': i_token
}

r = requests.request('GET', f"{i_vault}/v1/{i_engine}/data/{i_secret}", headers=headers)

data = json.loads(r.text)['data']['data']
data = remove_special(data)

if i_resource_type == 'configmap' or i_resource_type == 'config-map':
    apply_k8s_configmap(i_resource_name, data, i_namespace, i_kubeconfig)
else:
    apply_k8s_secret(i_resource_name, base64encode_collection(data), i_namespace, i_kubeconfig)