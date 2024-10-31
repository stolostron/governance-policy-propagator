[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Governance Policy Propagator [![KinD tests](https://github.com/stolostron/governance-policy-propagator/actions/workflows/kind.yml/badge.svg?branch=main&event=push)](https://github.com/stolostron/governance-policy-propagator/actions/workflows/kind.yml)

## Description

The governance policy propagator is a controller that watches `Policies`, `PlacementBindings`, and `PlacementRules`. It manages replicated Policies in cluster namespaces based on the PlacementBindings and PlacementRules, and it updates the status on Policies to show aggregated cluster compliance results. This controller is a part of the [governance-policy-framework](https://github.com/stolostron/governance-policy-framework).

The operator watches for changes to trigger a reconcile:

1. Changes to Policies in non-cluster namespaces trigger a self reconcile.
2. Changes to Policies in cluster namespaces trigger a root Policy reconcile.
2. Changes to PlacementBindings trigger reconciles on the subject Policies. 
3. Changes to PlacementRules trigger reconciles on subject Policies.

Every reconcile does the following:

1. Creates/updates/deletes replicated policies in cluster namespaces based on PlacementBinding/PlacementRule results.
2. Creates/updates/deletes the policy status to show aggregated cluster compliance results.

## Getting started

Go to the
[Contributing guide](https://github.com/open-cluster-management-io/community/blob/main/sig-policy/contribution-guidelines.md)
to learn how to get involved.

Check the [Security guide](SECURITY.md) if you need to report a security issue.

### Changes to the deploy YAML files

The YAML files in the deploy directory are autogenerated by Kubebuilder and Kustomize. After code
changes that affect the YAML files, the YAML files can be regenerated with
`make generate-operator-yaml`.

### Build and deploy locally
You will need [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) installed.

1. Create the Kind cluster
   ```bash
   make kind-bootstrap-cluster-dev
   ```
2. Start the propagator:
   - Run in a pod on the cluster:
     ```bash
     make build-images
     make kind-deploy-controller-dev
     ```
   - Run locally:
     ```bash
     make run
     ```

### Running tests
```
make test-dependencies
make test

make e2e-dependencies
make e2e-test
```

### How to run webhook locally

```bash
--enable-webhooks=true
```

> **Limit**: If you want to run the webhook locally, you need to generate certificates and place them in `/tmp/k8s-webhook-server/serving-certs/tls.{crt,key}`. If you’re not running a local API server, you’ll also need to figure out how to proxy traffic from the remote cluster to your local webhook server. For this reason, Kubebuilder generally recommends disabling webhooks when doing your local code-run-test cycle. To disable it, please supply the `--enable-webhooks=false` argument when running the controller.
> For more information, visit https://book.kubebuilder.io/cronjob-tutorial/running.html

### Clean up
```
make kind-delete-cluster
```

### Updating Deployment resources
Some of the deployment resources are generated by kubebuilder - the crds are generated into `./deploy/crds` and the rbac details from kubebuilder comments are compiled into `./deploy/rbac/role.yaml`.  Other details are managed independently - in particular, the details in `./deploy/manager/manager.yaml`. When any of those details need to be changed, the main deployment yaml `./deploy/operator.yaml` must be regenerated through the `make generate-operator-yaml` target. The `./deploy/operator.yaml` SHOULD NOT be manually updated.

## References

- The `governance-policy-propagator` is part of the `open-cluster-management` community. For more information, visit: [open-cluster-management.io](https://open-cluster-management.io).

<!---
Date: 07/18/2024
-->
