#!/usr/bin/env bash

# (C) Copyright NuoDB, Inc. 2019-2021  All Rights Reserved
# This file is licensed under the BSD 3-Clause License.
# See https://github.com/nuodb/nuodb-helm-charts/blob/master/LICENSE

# exit when any command fails
set -ex

# Download kubectl, which is a requirement for using minikube.
curl -Lo kubectl https://dl.k8s.io/release/v"${KUBERNETES_VERSION}"/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Download Helm
wget https://get.helm.sh/helm-"${HELM_VERSION}"-linux-amd64.tar.gz -O /tmp/helm.tar.gz
tar xzf /tmp/helm.tar.gz -C /tmp --strip-components=1 && chmod +x /tmp/helm && sudo mv /tmp/helm /usr/local/bin

mkdir -p $TEST_RESULTS # create the test results directory

if [[ "$REQUIRES_MINIKUBE" == "true" ]]; then
  sudo apt-get update

  # libncurses5 is needed for nuosql
  sudo apt-get install -y libncurses5 libncursesw5

  # Download minikube.
  curl -Lo minikube https://github.com/kubernetes/minikube/releases/v"${MINIKUBE_VERSION}"/download/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/

  # start minikube
  minikube start --vm-driver=docker --kubernetes-version=v"${KUBERNETES_VERSION}" --cpus=max --memory=max
  minikube status
  kubectl cluster-info

  # Inject database limits as the default CPUs limit of 8 is too big and causes
  # problems with minikube docker driver
  cat <<EOF > valuesInject.yaml
database:
  te:
    resources:
      limits:
        cpu: 2
        memory: 2Gi
  sm:
    resources:
      limits:
        cpu: 2
        memory: 2Gi
EOF

  # Configure DNS entries for Ingress testing
  ip=$(minikube ip)
  echo "$ip api.nuodb.local sql.nuodb.local demo.nuodb.local" | sudo tee -a /etc/hosts

  # Start 'minikube tunnel' so that services with type LoadBalancer are correctly
  # provisioned and routes to the minikube IP are created;
  # see https://minikube.sigs.k8s.io/docs/handbook/accessing/#using-minikube-tunnel
  nohup minikube tunnel > ${TEST_RESULTS}/minikube_tunnel.log 2>&1 &
  echo "echo \"MINIKUBE_PROCS: <\$(ps aux | grep minikube)>\"" >> $HOME/.nuodbrc

  # In some tests (specifically TestKubernetesTLSRotation), we observe incorrect DNS resolution
  # after pods have been re-created which causes problems with inter pod communication.
  # Set CoreDNS TTL to 0 so that DNS entries are not cached.
  kubectl get cm coredns -n kube-system -o yaml | sed -e 's/ttl [0-9]*$/ttl 0/' | kubectl apply -n kube-system -f -
  kubectl delete pods -l k8s-app=kube-dns -n kube-system

  helm version
  kubectl version

  # get the helm repo for upgrade testing
  helm repo add nuodb https://nuodb.github.io/nuodb-helm-charts
  helm repo add nuodb-incubator https://nuodb.github.io/nuodb-helm-charts/incubator

  # get HC Vault for testing
  helm repo add hashicorp https://helm.releases.hashicorp.com

  # get HAProxy for Ingress testing
  helm repo add haproxytech https://haproxytech.github.io/helm-charts

  # enable volume snapshots and CSI
  minikube addons enable volumesnapshots
  minikube addons enable csi-hostpath-driver

elif [[ "$REQUIRES_MINISHIFT" == "true" ]]; then
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
  helm repo add nuodb https://nuodb.github.io/nuodb-helm-charts
  helm repo add nuodb-incubator https://nuodb.github.io/nuodb-helm-charts/incubator
elif [[ "$REQUIRES_AKS" == "true" ]]; then
  if ! command -v az &> /dev/null; then
    # Install Azure cli if not already available
    curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
  fi
  # Azure login will be done by CircleCI automatically, however, let's keep the
  # code in case we need to migrate to other CI system. The AZURE_SP,
  # AZURE_SP_PASSWORD, AZURE_SP_TENANT, AZURE_AKS_RESOURCE_GROUP and
  # CLUSTER_CONTEXTS environment variables should be set before hand.
  if ! az account show > /dev/null 2>&1 ; then
    az login --service-principal -u "${AZURE_SP}" -p "${AZURE_SP_PASSWORD}" --tenant "${AZURE_SP_TENANT}"
  fi
  for cluster in ${CLUSTER_CONTEXTS}; do
    az aks get-credentials --name "$cluster" -g "${AZURE_AKS_RESOURCE_GROUP}"
  done
elif [[ "$REQUIRES_AWS" == "true" ]]; then
  for cluster in ${CLUSTER_CONTEXTS}; do
    if [[ "${cluster}" = "yin" ]]; then
      aws eks --region "${AWS_DEFAULT_REGION}" update-kubeconfig --name "${EKS_CLUSTER_1}" --alias "$cluster"
    else
      aws eks --region "${AWS_DEFAULT_REGION}" update-kubeconfig --name "${EKS_CLUSTER_2}" --alias "$cluster"
    fi
  done
else
  echo "Skipping installation steps"
fi

curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v"${GOTESTSUM_VERSION}"/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz" | sudo tar -xz -C /usr/local/bin gotestsum

# Install NuoDB client on the build host
curl -sSL "https://github.com/nuodb/nuodb-client/releases/download/v${NUODBCLIENT_VERSION}/nuodb-client-${NUODBCLIENT_VERSION}.lin-x64.tar.gz" | sudo tar -xz -C $HOME
echo "export PATH=${HOME}/nuodb-client-${NUODBCLIENT_VERSION}.lin-x64/bin:\$PATH" >> $HOME/.nuodbrc
echo "echo PATH: \"\$PATH\"" >> $HOME/.nuodbrc
