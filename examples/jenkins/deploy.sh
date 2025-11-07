#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Define variables
HELM_REPO_NAME="jenkinsci"
HELM_REPO_URL="https://charts.jenkins.io"
HELM_CHART_NAME="jenkins"
RELEASE_NAME="jenkins-dev"
NAMESPACE="jenkins"

# Add the Jenkins Helm repository
echo "Adding Jenkins Helm repository..."
helm repo add $HELM_REPO_NAME $HELM_REPO_URL

# Update your local Helm chart repository cache
echo "Updating Helm repositories..."
helm repo update

# Create the namespace if it doesn't exist
echo "Creating namespace '$NAMESPACE' if it doesn't exist..."
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Create PersistentVolume and PersistentVolumeClaim for Jenkins data
# echo "Creating PersistentVolume and PersistentVolumeClaim for Jenkins data..."
# kubectl apply -f jenkins-pvc.yaml -n $NAMESPACE

# Deploy the latest Jenkins chart with custom values
echo "Deploying Jenkins to namespace '$NAMESPACE'..."
helm upgrade --install $RELEASE_NAME $HELM_REPO_NAME/$HELM_CHART_NAME \
  --namespace $NAMESPACE \
  --values values.yaml

# Wait for the deployment to be ready
echo "Waiting for Jenkins to be ready..."
kubectl wait --namespace $NAMESPACE \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=jenkins-controller \
  --timeout=300s

echo "Jenkins has been deployed."
echo "You can access it by port-forwarding the service:"
echo "kubectl --namespace $NAMESPACE port-forward svc/$RELEASE_NAME 8080:8080"

echo "To get the admin password, run:"
echo "kubectl exec --namespace $NAMESPACE -it svc/$RELEASE_NAME -c jenkins -- /bin/cat /run/secrets/additional/chart-admin-password && echo"
