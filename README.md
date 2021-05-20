[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Governance Policy Propagator [![KinD tests](https://github.com/open-cluster-management/governance-policy-propagator/actions/workflows/kind.yml/badge.svg?branch=main&event=push)](https://github.com/open-cluster-management/governance-policy-propagator/actions/workflows/kind.yml)

## Description

The governance policy propagator is a controller that watches `Policies`, `PlacementBindings`, and `PlacementRules`. It manages replicated Policies in cluster namespaces based on the PlacementBindings and PlacementRules, and it updates the status on Policies to show aggregated cluster compliance results. The controller can be run as a stand-alone program or as an integrated part of governing risk with the Open Cluster Management project.

Go to the [Contributing guide](CONTRIBUTING.md) to learn how to get involved.

## Geting started 

Check the [Security guide](SECURITY.md) if you need to report a security issue.

### Build and deploy locally
You will need [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) installed.

```bash
make kind-bootstrap-cluster-dev
make build-images
make kind-deploy-controller-dev
```
### Running tests
```
make test-dependencies
make test

make e2e-dependencies
make e2e-test
```

### Clean up
```
make kind-delete-cluster
```

## References

- The `governance-policy-propagator` is part of the `open-cluster-management` community. For more information, visit: [open-cluster-management.io](https://open-cluster-management.io).
