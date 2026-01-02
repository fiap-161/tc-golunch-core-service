# 🔐 Arquitetura de Autenticação - GoLunch (Serverless + Microservices)

## 📐 **Visão Geral da Arquitetura Atual**

### ✅ **Implementação Real (Janeiro 2026):**
- **AUTENTICAÇÃO**: Serverless (AWS Lambda + API Gateway + Cognito)
- **CORE SERVICES**: 3 Microserviços (Core, Payment, Operation)
- **COMUNICAÇÃO**: Service-to-service authentication via API keys

### 🏗️ **Arquitetura Atual:**

```
┌─────────────────────────────────────────────────────────────────┐
│                🚀 SERVERLESS AUTH (AWS)                        │
│                 tc-golunch-serverless/                          │
│  ┌─────────────────────┐  ┌─────────────────────────────────┐   │
│  │   🌐 API Gateway    │  │       ⚡ Lambda Functions       │   │
│  │                     │  │                                 │   │
│  │ POST /auth/register │◄─┤ RegisterUser                    │   │
│  │ POST /auth/login    │◄─┤ LoginUser                       │   │
│  │ GET /auth/anonymous │◄─┤ AnonymousLogin                  │   │
│  │ POST /admin/register│◄─┤ AdminRegister                   │   │
│  │ POST /admin/login   │◄─┤ AdminLogin                      │   │
│  │ POST /validate-service│◄┤ ServiceAuth (Service-to-Service)│   │
│  └─────────────────────┘  └─────────────────────────────────┘   │
│           │                        │                             │
│           │                ┌───────▼──────┐                     │
│           │                │  🔐 Cognito  │                     │
│           │                │  User Pool   │                     │
│           │                └──────────────┘                     │
└───────────┼────────────────────────────────────────────────────┘
            │
    ┌───────▼─────────────────────────────────────────┐
    │           🎯 MICROSERVICES (Local/EKS)          │
    │                                                 │
    │  ┌──────────────┐ ┌─────────────┐ ┌─────────┐  │
    │  │ Core Service │ │Payment Srvc │ │Operation│  │
    │  │   (8081)     │ │   (8082)    │ │Service  │  │
    │  │              │ │             │ │ (8083)  │  │
    │  │• Products    │ │• MercadoPago│ │• Kitchen│  │
    │  │• Orders      │ │• Webhooks   │ │• Status │  │
    │  │• Customers   │ │• Payments   │ │• Admin  │  │
    │  └──────────────┘ └─────────────┘ └─────────┘  │
    │        │               │               │       │
    │        └──── Service-to-Service Auth ──┘       │
    │           (API Keys via headers)                │
    └─────────────────────────────────────────────────┘
```

## 🎯 **Responsabilidades por Serviço**

### � **Serverless Auth (AWS Lambda)**
**RESPONSABILIDADES:**
- ✅ **Customer Authentication**: Register, Login, Anonymous tokens
- ✅ **Admin Authentication**: Register, Login, Token validation  
- ✅ **JWT Generation/Validation**: Via Amazon Cognito + Lambda
- ✅ **Service-to-Service Auth**: Validação de API keys entre microservices
- ✅ **Scalable & Cost-effective**: Paga apenas quando usar

**ENDPOINTS (API Gateway):**
```bash
# Customer Authentication
POST   /auth/register          # Register new customer → RegisterUser Lambda
POST   /auth/login             # Customer login → LoginUser Lambda
GET    /auth/anonymous         # Generate anonymous token → AnonymousLogin Lambda

# Admin Authentication  
POST   /admin/register         # Register new admin → AdminRegister Lambda
POST   /admin/login            # Admin login → AdminLogin Lambda

# Service-to-Service Authentication
POST   /validate-service       # Validate API keys → ServiceAuth Lambda
```

### 🔑 **Core Service (8081)**
**RESPONSABILIDADES:**
- ✅ **Product Management**: CRUD operations 
- ✅ **Order Management**: CRUD operations
- ✅ **Customer Management**: Basic customer data (não auth)
- ❌ **REMOVED**: Authentication logic (moved to serverless)

**ENDPOINTS:**
```bash
# Product Management
GET    /product/categories     # List categories
GET    /product               # List products

# Order Management  
POST   /order                 # Create order
GET    /order                 # List orders
PUT    /order/:id             # Update order

# Customer Management (data only)
GET    /customer/identify/:cpf # Get customer data
POST   /customer/register     # Store customer data
```

### 💳 **Payment Service (8082)**
**RESPONSABILIDADES:**
- ✅ **Payment Processing**: MercadoPago integration
- ✅ **Webhook Handling**: Payment status updates
- ✅ **Token Validation**: Via HTTP calls to Auth Service

### 🏭 **Operation Service (8083)**  
**RESPONSABILIDADES:**
- ✅ **Production Management**: Kitchen operations
- ✅ **Order Status Updates**: Production workflow
- ✅ **Admin Panel**: Dashboard para cozinha
- ✅ **Token Validation**: Via Service-to-Service auth
- ❌ **REMOVED**: Admin authentication (now serverless)

## 🔑 **Fluxo de Autenticação Atual**

### **👤 Customer Flow (via Serverless):**
```bash
# 1. Register via API Gateway → Lambda
POST https://api-gateway-url.amazonaws.com/auth/register
{
  "name": "João Silva",
  "email": "joao@example.com", 
  "cpf": "12345678901"
}

# 2. Login via API Gateway → Lambda
POST https://api-gateway-url.amazonaws.com/auth/login
{
  "email": "joao@example.com",
  "password": "password123"
}
# Returns: {"token": "eyJhbGci...", "user_type": "customer"}

# 3. Use token for orders (Core Service)
Authorization: Bearer <customer_token>
POST http://localhost:8081/order
```

### **👑 Admin Flow (via Serverless):**
```bash
# 1. Register Admin via API Gateway → Lambda
POST https://api-gateway-url.amazonaws.com/admin/register
{
  "email": "admin@golunch.com",
  "password": "admin123456"
}

# 2. Login Admin via API Gateway → Lambda
POST https://api-gateway-url.amazonaws.com/admin/login  
{
  "email": "admin@golunch.com", 
  "password": "admin123456"
}
# Returns: {"token": "eyJhbGci...", "user_type": "admin"}

# 3. Use token for admin operations (Operation Service)
Authorization: Bearer <admin_token>
GET http://localhost:8083/admin/orders
```

### **🔍 Service-to-Service Authentication:**
```bash
# Core Service → Payment Service
X-Service-Name: core-service
X-Service-Key: core-service-secure-api-key-2024
POST http://localhost:8082/payment

# Validation via Serverless (if configured)
POST https://api-gateway-url.amazonaws.com/validate-service
{
  "serviceName": "core-service",
  "apiKey": "core-service-secure-api-key-2024"
}
# Returns: {"success": true, "serviceName": "core-service"}
```

## 💾 **Database Schema**

### **Serverless (Amazon Cognito User Pool):**
```json
// Cognito manages users automatically
{
  "customer": {
    "sub": "uuid",
    "email": "customer@example.com",
    "custom:cpf": "12345678901",
    "custom:name": "João Silva",
    "custom:user_type": "customer"
  },
  "admin": {
    "sub": "uuid", 
    "email": "admin@golunch.com",
    "custom:user_type": "admin"
  }
}
```

### **Core Service Database (port 5433):**
```sql
-- Customer data (não auth)
CREATE TABLE customer_daos (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    email VARCHAR,
    cpf VARCHAR UNIQUE,
    is_anonymous BOOLEAN DEFAULT false,
    cognito_sub VARCHAR, -- Reference to Cognito user
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Products table
CREATE TABLE product_daos (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    category VARCHAR,
    price DECIMAL,
    description TEXT,
    image_url VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Orders table  
CREATE TABLE order_daos (
    id VARCHAR PRIMARY KEY,
    customer_id VARCHAR,
    status VARCHAR,
    total DECIMAL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## 🧪 **Testing with Bruno Collections**

### **Updated Bruno Structure:**
```
fiap161-tc-collections/
├── core-service/              # Core business logic
│   ├── CREATE CUSTOMER.bru    # Customer data (not auth)
│   ├── LIST PRODUCTS.bru
│   ├── CREATE ORDER.bru
│   └── ...
├── payment-service/           # Payment processing
│   ├── CREATE PAYMENT.bru
│   └── ...
├── operation-service/         # Kitchen operations
│   ├── LIST ORDERS.bru
│   └── ...
└── microservices-auth/        # Service-to-service auth
    ├── SERVICE AUTH - Health Check.bru
    ├── SERVICE AUTH - Create Payment.bru
    └── ...
```

### **Serverless Testing (separado):**
```
tc-golunch-serverless/
├── test-auth.sh               # Test customer/admin auth
├── test-service-auth.sh       # Test service-to-service
└── auth/                      # Lambda function code
    ├── register.js
    ├── login.js
    ├── admin-register.js
    ├── admin-login.js
    └── service-auth.js
```

## ✅ **Benefits of Serverless Architecture**

### **🎯 Eliminated Problems:**
1. **❌ Infrastructure Management**: AWS manages servers
2. **❌ Scaling Issues**: Lambda scales automatically  
3. **❌ High Costs**: Pay only for execution time
4. **❌ Complex Deployments**: Simple zip upload

### **✅ Achieved Benefits:**
1. **� Auto-scaling**: Handle traffic spikes automatically
2. **💰 Cost Effective**: No idle server costs
3. **� Managed Security**: Cognito handles user management
4. **🧹 Clean Separation**: Auth completely independent
5. **🌎 Global Distribution**: Edge locations via API Gateway
6. **📈 Monitoring**: CloudWatch logs and metrics built-in

### **⚡ Performance Benefits:**
- Zero cold start for Cognito operations
- API Gateway caching for repeated calls  
- Lambda concurrent execution for high load
- No database connections to manage

## 🚀 **Migration Guide & Current State**

### **Development vs Production:**
```bash
# Local Development (Microservices only)
export SERVICE_AUTH_API_URL=""  # Empty = local validation
./setup-microservices.sh
# Services communicate via direct HTTP + API keys

# Production (Serverless + EKS)
export SERVICE_AUTH_API_URL="https://api-gateway-url.amazonaws.com/validate-service"
# Services can optionally validate via Lambda
# Fallback to local validation if Lambda unavailable
```

### **Deployment Strategy:**
```bash
# 1. Deploy Serverless Auth (AWS)
cd tc-golunch-serverless/
terraform apply

# 2. Deploy Microservices (EKS)  
cd tc-golunch-infra/
terraform apply

# 3. Configure service URLs
export COGNITO_AUTH_URL="https://api-gateway-url.amazonaws.com"
export SERVICE_AUTH_API_URL="https://api-gateway-url.amazonaws.com/validate-service"
```

### **For Frontend/API Clients:**
```bash
# Authentication Endpoints (Serverless)
POST https://api-gateway-url.amazonaws.com/auth/register
POST https://api-gateway-url.amazonaws.com/auth/login
POST https://api-gateway-url.amazonaws.com/admin/login

# Business Logic Endpoints (Microservices)
GET http://core-service-url/product
POST http://core-service-url/order
GET http://operation-service-url/admin/orders
```

---

## 📊 **Architecture Comparison**

| Aspect | Serverless Auth | Microservices Auth | Monolith Auth |
|--------|----------------|-------------------|---------------|
| **Scaling** | Automatic | Manual K8s | Manual servers |
| **Cost** | Pay-per-use | Fixed costs | Fixed costs |
| **Maintenance** | AWS managed | Self-managed | Self-managed |
| **Complexity** | Low | Medium | High |
| **Cold Start** | ~100ms | Always warm | Always warm |
| **Availability** | 99.95% SLA | Depends on infra | Depends on infra |

**🎯 Result**: Hybrid architecture with serverless auth + microservices business logic provides best of both worlds.

---

## 🔗 **Additional Resources**

- **📁 tc-golunch-serverless/**: Complete serverless implementation
- **📋 MICROSERVICES_DOCUMENTATION.md**: Full microservices guide  
- **🧪 SERVICE_TO_SERVICE_AUTH_STATUS.md**: Service auth implementation
- **🚀 TESTING_LOCAL.md**: Local development setup