# Crypto Kubernetes Operator

A Kubernetes Operator built with `client-go` to automate the deployment of the Crypto Aggregator stack.

## Features
- **Namespace Watcher:** Monitors for namespaces with the `crypto-` prefix.
- **Automated Provisioning:** Deploys a PostgreSQL instance (Deployment + Service) and the Aggregator application upon namespace creation.
- **Reconciliation:** Detects deletion of the aggregator deployment and automatically recreates it.
- **Dependency Injection:** Programmatically injects database credentials and hostnames into the application environment.

## Prerequisites
- Go 1.25
- Kind / Local K8s cluster
- Kubectl

## Getting Started

1. **Prepare Cluster:**
   ```bash
   kind create cluster --name crypto-dev
   docker build -t crypto-aggregator:v2 ../crypto-aggregator
   kind load docker-image crypto-aggregator:v2 --name crypto-dev
   ```

2. **Run Operator:**
   ```bash
   go run main.go
   ```

3. **Deploy Stack:**
   ```bash
   kubectl create ns crypto-test
   ```

## Verification
```bash
# Check resources
kubectl get all -n crypto-test

# Check app logs
kubectl logs -n crypto-test -l app=aggregator
```