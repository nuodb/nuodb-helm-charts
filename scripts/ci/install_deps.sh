#!/usr/bin/env bash

# Download kubectl, which is a requirement for using minikube.
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v"${KUBERNETES_VERSION}"/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
# Download minikube.
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v"${MINIKUBE_VERSION}"/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
mkdir -p "$HOME"/.kube "$HOME"/.minikube
touch "$KUBECONFIG"

# start minikube
sudo minikube start --vm-driver=none --kubernetes-version=v"${KUBERNETES_VERSION}" --memory=8000 --cpus=4
sudo chown -R travis: /home/travis/.minikube/
kubectl cluster-info

# install Helm and Tiller
wget http://storage.googleapis.com/kubernetes-helm/helm-"${HELM_VERSION}"-linux-amd64.tar.gz -O /tmp/helm.tar.gz
tar xzf /tmp/helm.tar.gz -C /tmp --strip-components=1 && chmod +x /tmp/helm && sudo mv /tmp/helm /usr/local/bin
helm init