#!/usr/bin/env bash

# exit when any command fails
set -e

# Download kubectl, which is a requirement for using minikube.
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v"${KUBERNETES_VERSION}"/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Download Helm and Tiller
wget https://get.helm.sh/helm-"${HELM_VERSION}"-linux-amd64.tar.gz -O /tmp/helm.tar.gz
tar xzf /tmp/helm.tar.gz -C /tmp --strip-components=1 && chmod +x /tmp/helm && sudo mv /tmp/helm /usr/local/bin

if [[ -n "$REQUIRES_MINIKUBE" ]]; then
    # Download minikube.
  curl -Lo minikube https://storage.googleapis.com/minikube/releases/v"${MINIKUBE_VERSION}"/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
  mkdir -p "$HOME"/.kube "$HOME"/.minikube
  touch "$KUBECONFIG"

  # start minikube
  sudo minikube start --vm-driver=none --kubernetes-version=v"${KUBERNETES_VERSION}" --memory=8000 --cpus=4
  sudo chown -R travis: /home/travis/.minikube/
  kubectl cluster-info

  # In some tests (specifically TestKubernetesTLSRotation), we observe incorrect DNS resolution 
  # after pods have been re-created which causes problems with inter pod communicataion.
  # Set CoreDNS TTL to 0 so that DNS entries are not cached. 
  kubectl get cm coredns -n kube-system -o yaml | sed -e 's/ttl [0-9]*$/ttl 0/' | kubectl apply -n kube-system -f -
  kubectl delete pods -l k8s-app=kube-dns -n kube-system

  helm version

  kubectl version

  # get the helm repo for upgrade testing
  helm repo add nuodb https://storage.googleapis.com/nuodb-charts

elif [[ -n "$REQUIRES_MINISHIFT" ]]; then
  wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz -O /tmp/oc.tar.gz
  tar xzf /tmp/oc.tar.gz -C /tmp --strip-components=1 && chmod +x /tmp/oc && sudo mv /tmp/oc /usr/local/bin

  oc version

  sudo apt install libvirt-bin qemu-kvm
  sudo usermod -a -G libvirtd "$(whoami)"

  curl -L https://github.com/dhiltgen/docker-machine-kvm/releases/download/v0.10.0/docker-machine-driver-kvm-ubuntu14.04 -o /tmp/docker-machine-driver-kvm
  chmod +x /tmp/docker-machine-driver-kvm && sudo mv /tmp/docker-machine-driver-kvm /usr/local/bin

  wget https://github.com/minishift/minishift/releases/download/v"${MINISHIFT_VERSION}"/minishift-"${MINISHIFT_VERSION}"-linux-amd64.tgz -O /tmp/minishift.tar.gz
  tar xzf /tmp/minishift.tar.gz -C /tmp --strip-components=1 && chmod +x /tmp/minishift && sudo mv /tmp/minishift /usr/local/bin

  sudo minishift start --openshift-version=v"${OPENSHIFT_VERSION}" --memory=8000 --cpus=4

  oc login -u system:admin
  oc status

  kubectl cluster-info

  kubectl -n kube-system create serviceaccount tiller-system
  kubectl create clusterrolebinding tiller-system --clusterrole cluster-admin --serviceaccount=kube-system:tiller-system

  helm version

  kubectl version

  # disable THP to match minikube
  kubectl create namespace nuodb
  kubectl -n nuodb create serviceaccount nuodb
  oc apply -f deploy/nuodb-scc.yaml -n nuodb
  oc adm policy add-scc-to-user nuodb-scc system:serviceaccount:nuodb:nuodb -n nuodb
  oc adm policy add-scc-to-user nuodb-scc system:serviceaccount:nuodb:default -n nuodb
  helm install stable/transparent-hugepage/ --namespace nuodb

  # get the helm repo for upgrade testing
  helm repo add nuodb https://storage.googleapis.com/nuodb-charts
else
  echo "Skipping installation steps"
fi