#!/bin/sh

printenv

set -e

# Extract the base64 encoded config data and write this to the KUBECONFIG
mkdir -p ~/.kube
echo $INPUT_KUBECONFIG | base64 -d > ~/.kube/config

echo 'current kubectl context: '
kubectl config current-context

python -u /usr/local/bin/main.py