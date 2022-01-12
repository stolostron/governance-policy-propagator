#!/bin/bash
set -e

_OS=$(go env GOOS)
_ARCH=$(go env GOARCH)

if ! which kubectl > /dev/null; then
    echo "installing kubectl"
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
fi
if ! which kind > /dev/null; then
    echo "installing kind"
    curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.8.1/kind-$(uname)-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
fi
echo "Installing ginkgo ..."
go get github.com/onsi/ginkgo/ginkgo@v1.14.1
go get github.com/onsi/gomega/...@v1.10.1

# echo "Install Kubebuilder components for test framework usage!"
# download kubebuilder and extract it to tmp
# curl -L https://go.kubebuilder.io/dl/2.2.0/"${_OS}"/"${_ARCH}" | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
# sudo mv /tmp/kubebuilder_2.2.0_"${_OS}"_"${_ARCH}" /usr/local/kubebuilder
# export PATH=$PATH:/usr/local/kubebuilder/bin

# Run unit test
export IMAGE_NAME_AND_VERSION=${1}
# make test
# run unit test
go test `go list ./... | grep -v test/e2e` > report.json
# run e2e test
make build-instrumented
make kind-bootstrap-cluster-dev
make run-instrumented
make e2e-test
