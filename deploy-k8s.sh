#!/bin/bash

# Deploy Core Service to Kubernetes
# Usage: ./deploy-k8s.sh [namespace]

NAMESPACE=${1:-golunch}

echo "🚀 Deploying Core Service to namespace: ${NAMESPACE}"
echo "💰 Cost: $0 (using PostgreSQL StatefulSet)"

# Create namespace if it doesn't exist
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

echo "🗄️ Deploying PostgreSQL..."
kubectl apply -f k8s/postgres-statefulset.yaml -n ${NAMESPACE}

echo "⏳ Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres-core -n ${NAMESPACE} --timeout=300s

echo "� Creating PVC for uploads..."
kubectl apply -f k8s/core-service-uploads-pvc.yaml -n ${NAMESPACE}

echo "�📦 Applying ConfigMap..."
kubectl apply -f k8s/core-service-configmap.yaml -n ${NAMESPACE}

echo "🔐 Applying Secrets..."
kubectl apply -f k8s/core-service-secrets.yaml -n ${NAMESPACE}

echo "🚀 Applying Deployment..."
kubectl apply -f k8s/core-service-deployment.yaml -n ${NAMESPACE}

echo "🌐 Applying Service..."
kubectl apply -f k8s/core-service-service.yaml -n ${NAMESPACE}

echo "📈 Applying HPA..."
kubectl apply -f k8s/core-service-hpa.yaml -n ${NAMESPACE}

# Wait for deployment to be ready
echo "⏳ Waiting for Core Service to be ready..."
kubectl rollout status deployment/core-service -n ${NAMESPACE} --timeout=300s

# Show deployment status
echo ""
echo "✅ Deployment Status:"
kubectl get pods -l app=core-service -n ${NAMESPACE}
kubectl get pods -l app=postgres-core -n ${NAMESPACE}
kubectl get svc -n ${NAMESPACE}

echo ""
echo "🎉 Core Service deployed successfully!"
echo ""
echo "📊 Next Steps:"
echo "  • Test: kubectl port-forward svc/core-service 8081:8081 -n ${NAMESPACE}"
echo "  • Check: curl http://localhost:8081/ping"
echo "  • Logs: kubectl logs -f deployment/core-service -n ${NAMESPACE}"
echo "  • DB Access: kubectl port-forward svc/postgres-core 5432:5432 -n ${NAMESPACE}"