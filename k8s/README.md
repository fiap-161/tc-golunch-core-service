# Kubernetes Deployment - Core Service

Deploy simples do **GoLunch Core Service** com PostgreSQL StatefulSet.

## 💰 **CUSTO ZERO**
- ✅ PostgreSQL StatefulSet no cluster
- ✅ Sem custos AWS adicionais
- ✅ Funcional para desenvolvimento e testes

---

## 📋 **Arquivos**

```
k8s/
├── postgres-statefulset.yaml      # PostgreSQL + Service
├── core-service-deployment.yaml   # Core Service app
├── core-service-service.yaml      # Services e LoadBalancer  
├── core-service-configmap.yaml    # Configurações
├── core-service-secrets.yaml      # Credenciais
├── core-service-hpa.yaml          # Auto-scaling
└── README.md
```

---

## 🚀 **Deploy Rápido**

```bash
# Deploy completo (PostgreSQL + Core Service)
./deploy-k8s.sh

# Ou especificar namespace
./deploy-k8s.sh meu-namespace
```

---

## 🧪 **Testar**

```bash
# Port-forward para testar
kubectl port-forward svc/core-service 8081:8081 -n golunch

# Test API
curl http://localhost:8081/ping

# Logs
kubectl logs -f deployment/core-service -n golunch
```

---

## 🗄️ **Acesso ao Database**

```bash
# Port-forward PostgreSQL
kubectl port-forward svc/postgres-core 5432:5432 -n golunch

# Conectar via psql
psql -h localhost -U golunch_core_user -d golunch_core
# Password: golunch_core_password
```

---

## 📊 **Recursos**

- **PostgreSQL**: 1 replica, 10GB storage
- **Core Service**: 3 replicas, auto-scaling 3-10
- **Memory**: 128Mi-256Mi por pod
- **CPU**: 200m-500m por pod

---

## 🔧 **Comandos Úteis**

```bash
# Status geral
kubectl get all -n golunch

# Restart deployment
kubectl rollout restart deployment/core-service -n golunch

# Backup database
kubectl exec postgres-core-0 -n golunch -- \
  pg_dump -U golunch_core_user golunch_core > backup.sql

# Cleanup
kubectl delete namespace golunch
```

---
