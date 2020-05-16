#!/bin/bash
set -e
echo "UNIT TESTS GO HERE!"

echo "Install Kubebuilder components for test framework usage!"

_OS=$(go env GOOS)
_ARCH=$(go env GOARCH)

if ! which kind > /dev/null; then
    echo "installing kind"
    curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.8.1/kind-$(uname)-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
fi

# download kubebuilder and extract it to tmp
# curl -L https://go.kubebuilder.io/dl/2.2.0/"${_OS}"/"${_ARCH}" | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
# sudo mv /tmp/kubebuilder_2.2.0_"${_OS}"_"${_ARCH}" /usr/local/kubebuilder
# export PATH=$PATH:/usr/local/kubebuilder/bin

# Run unit test
export IMAGE_NAME_AND_VERSION=${1}
# make test
make build-instrumented-profile
make kind-bootstrap-cluster-dev
make run-instrumented-profile
make e2e-test
make stop-instrumented-profile
cat coverage.out
gosec -fmt sonarqube -out gosec.json -no-fail ./...
unset SONARQUBE_SCANNER_PARAMS
sonar-scanner --debug
