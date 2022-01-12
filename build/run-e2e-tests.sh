#!/bin/bash

set -e


CURR_FOLDER_PATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
KIND_KUBECONFIG="${CURR_FOLDER_PATH}/../kind_kubeconfig.yaml"
export KUBECONFIG=${KIND_KUBECONFIG}
export DOCKER_IMAGE_AND_TAG=${1}

if ! which kubectl > /dev/null; then
    echo "installing kubectl"
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
fi
if ! which kind > /dev/null; then
    echo "installing kind"
    curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.9.0/kind-$(uname)-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
fi
echo "Installing ginkgo ..."
go get github.com/onsi/ginkgo/ginkgo@v1.14.1
go get github.com/onsi/gomega/...@v1.10.1

make kind-create-cluster 

# setup kubeconfig
kind get kubeconfig --name test-hub > ${KIND_KUBECONFIG}

make install-crds 

make kind-deploy-controller 

echo "patch image"
kubectl patch deployment governance-policy-propagator -n governance -p "{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"governance-policy-propagator\",\"image\":\"${DOCKER_IMAGE_AND_TAG}\"}]}}}}"
kubectl rollout status -n governance deployment governance-policy-propagator --timeout=90s
sleep 10

make install-resources

make e2e-test

echo "delete cluster"
make kind-delete-cluster 

