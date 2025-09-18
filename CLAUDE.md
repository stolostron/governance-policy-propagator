# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

The governance-policy-propagator is a Kubernetes controller for Open Cluster Management that manages policy propagation across clusters. It watches `Policies`, `PlacementBindings`, and `PlacementRules` to create and manage replicated policies in cluster namespaces based on placement decisions.

## Common Development Commands

### Building and Testing
```bash
# Build the project
make build

# Run tests
make test-dependencies
make test

# Run end-to-end tests
make e2e-dependencies
make e2e-test

# Test coverage
make test-coverage

# Generate manifests and code
make generate
make manifests
```

### Local Development with Kind
```bash
# Create Kind cluster for development
make kind-bootstrap-cluster-dev

# Build and deploy controller in cluster
make build-images
make kind-deploy-controller-dev

# Run controller locally (for development)
make run

# Clean up Kind cluster
make kind-delete-cluster
```

### YAML Generation
```bash
# Regenerate deployment YAML (required after code changes affecting YAML)
make generate-operator-yaml
```

## Architecture Overview

### Core Controllers

1. **Root Policy Controller** (`controllers/propagator/rootpolicy_controller.go`)
   - Handles root policies in hub cluster
   - Manages placement and replication to cluster namespaces
   - Updates policy status with placement information

2. **Replicated Policy Controller** (`controllers/propagator/replicatedpolicy_controller.go`)
   - Manages policies in cluster namespaces
   - Handles template resolution and encryption
   - Processes policy updates from placement changes

3. **Policy Set Controller** (`controllers/policyset/`)
   - Manages PolicySet resources and their dependencies
   - Handles policy grouping and propagation

4. **Policy Automation Controller** (`controllers/automation/`)
   - Manages PolicyAutomation resources for Ansible integration
   - Handles automation triggers based on policy compliance

5. **Root Policy Status Controller** (`controllers/rootpolicystatus/`)
   - Aggregates compliance status from cluster policies back to root policies

6. **Encryption Keys Controller** (`controllers/encryptionkeys/`)
   - Manages policy template encryption key rotation
   - Default rotation period: 30 days

7. **Policy Metrics Controller** (`controllers/policymetrics/`)
   - Collects and reports policy compliance metrics

### Key Components

- **Propagator** (`controllers/propagator/propagation.go`): Core propagation logic
- **Template Resolution** (`controllers/propagator/template_utils.go`): Handles Go template processing in policies
- **Encryption** (`controllers/propagator/encryption.go`): Policy template encryption/decryption
- **Common Utilities** (`controllers/common/`): Shared controller utilities

### API Types

- **v1**: Main policy API (`api/v1/policy_types.go`, `api/v1/placementbinding_types.go`)
- **v1beta1**: PolicySet and PolicyAutomation APIs (`api/v1beta1/`)

## Configuration

### Environment Variables
- `WATCH_NAMESPACE`: Namespace(s) to watch (required)
- `DISABLE_REPORT_METRICS`: Set to "true" to disable metrics reporting
- `CONTROLLER_CONFIG_QPS`: Override default client QPS (200.0)
- `CONTROLLER_CONFIG_BURST`: Override default client burst (400)

### Command-line Flags
- `--enable-webhooks`: Enable policy validation webhook (default: true)
- `--disable-placementrule`: Disable PlacementRule watches
- `--encryption-key-rotation`: Key rotation days (default: 30)
- `--leader-elect`: Enable leader election (default: true)
- Various concurrency settings for different controllers

## Testing

- Unit tests: `make test`
- E2E tests in `test/e2e/` using Ginkgo framework
- Webhook-specific tests: `make e2e-test-webhook`
- PolicyAutomation tests: `make e2e-test-policyautomation`
- Coverage target: 71% minimum

## Development Notes

- Uses controller-runtime framework for Kubernetes controllers
- Supports both PlacementRules and Placements for policy targeting
- Template resolution supports Go templates with custom delimiters
- Policy encryption uses AES encryption for sensitive template data
- Webhook validation for policy resources when enabled
- Metrics exposed on port 8383, health probes on 8081

## Deployment

- Main deployment YAML: `deploy/operator.yaml` (auto-generated, do not edit manually)
- CRDs in `deploy/crds/`
- RBAC in `deploy/rbac/`
- Webhook configuration in `deploy/webhook.yaml`