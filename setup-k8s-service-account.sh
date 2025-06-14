#!/bin/bash

# Script to set up Kubernetes service account for GitHub Actions
set -e

echo "Setting up Kubernetes service account for GitHub Actions..."

# Apply the service account configuration
kubectl apply -f k8s-service-account.yaml

# Wait for the token to be created
echo "Waiting for service account token..."
sleep 5

# Get the token
TOKEN=$(kubectl get secret github-actions-deployer-token -o jsonpath='{.data.token}' | base64 -d)
CA_CERT=$(kubectl get secret github-actions-deployer-token -o jsonpath='{.data.ca\.crt}')

# Get cluster server URL
SERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')

# Create kubeconfig
cat << EOF > github-actions-kubeconfig.yaml
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: $CA_CERT
    server: $SERVER
  name: github-actions-cluster
contexts:
- context:
    cluster: github-actions-cluster
    user: github-actions-deployer
  name: github-actions-context
current-context: github-actions-context
users:
- name: github-actions-deployer
  user:
    token: $TOKEN
EOF

echo "✅ Service account created successfully!"
echo ""
echo "To use this with GitHub Actions:"
echo "1. Encode the kubeconfig:"
echo "   cat github-actions-kubeconfig.yaml | base64 -w 0"
echo ""
echo "2. Add the output as KUBECONFIG secret in GitHub repository settings"
echo ""
echo "3. The workflow will automatically use this configuration"

# Clean up
rm -f github-actions-kubeconfig.yaml

echo ""
echo "Service account setup completed!"
