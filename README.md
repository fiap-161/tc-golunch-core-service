# 🍔 GoLunch Core Service

Microsserviço central responsável pelas funcionalidades essenciais da lanchonete. Este serviço implementa a lógica de negócio para autenticação, gestão de pedidos, catálogo de produtos e clientes.

## 🎯 Responsabilidades

- **🔐 Autenticação**: Login de clientes e administradores com JWT
- **📋 Gerenciamento de Pedidos**: Criação, consulta e atualização de pedidos
- **📦 Catálogo de Produtos**: Listagem e consulta de produtos por categoria  
- **👥 Gestão de Clientes**: Identificação e cadastro de clientes
- **🔗 Relacionamento Pedido-Produto**: Associação de produtos aos pedidos
- **📊 Status de Pedidos**: Controle do fluxo de status dos pedidos
- **🛡️ Autorização**: Validação de permissões de admin

## 🏗️ Arquitetura

O serviço segue os princípios da **Arquitetura Hexagonal** com as seguintes camadas:

- **Entities**: Regras de negócio fundamentais
- **Use Cases**: Lógica de negócio específica
- **Gateways**: Interfaces para acesso a dados externos
- **Controllers**: Coordenação entre camadas
- **Handlers**: Gerenciamento de requisições HTTP
- **External/Infrastructure**: Implementações concretas (banco de dados)

## 🗄️ Banco de Dados

- **PostgreSQL**: Banco de dados principal
- **Tabelas**:
  - `customers`: Dados dos clientes
  - `products`: Catálogo de produtos
  - `orders`: Pedidos realizados
  - `product_orders`: Relacionamento pedido-produto

## 🚀 Endpoints Disponíveis

### 🔐 Autenticação
- `POST /admin/login` - Login de administrador
- `POST /admin/register` - Cadastro de administrador

### 👥 Clientes
- `GET /customer/identify/:cpf` - Identificar cliente por CPF
- `GET /customer/anonymous` - Login anônimo
- `POST /customer/register` - Cadastrar novo cliente

### 📦 Produtos
- `GET /product/categories` - Listar categorias de produtos
- `GET /product` - Listar produtos por categoria

### 📋 Pedidos
- `POST /order` - Criar novo pedido
- `GET /order` - Listar todos os pedidos
- `PUT /order/:id` - Atualizar pedido
- `GET /order/panel` - Painel de pedidos

### 🏥 Health Check
- `GET /ping` - Health check do serviço

## 🔧 Configuração Local

1. **Clone o repositório**
2. **Configure as variáveis de ambiente**:
   ```bash
   export DATABASE_URL="postgres://user:password@localhost:5432/golunch_core?sslmode=disable"
   export UPLOAD_DIR="./uploads"
   ```

3. **Execute o banco de dados**:
   ```bash
   docker-compose up -d postgres
   ```

4. **Execute a aplicação**:
   ```bash
   go run cmd/api/main.go
   ```

## 📋 Dependências

- **Go** 1.24.3
- **PostgreSQL** 16.3
- **Gin** - Framework web
- **GORM** - ORM para banco de dados
- **Swagger** - Documentação da API

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes BDD
go test -tags=bdd ./...
```

## 📊 Cobertura de Testes

- **Meta**: 80% de cobertura
- **BDD**: Implementado para cenários de criação de pedidos
- **Testes Unitários**: Todos os use cases e controllers

## 🐳 Docker

```bash
# Build da imagem
docker build -t tc-golunch-core-service .

# Executar container
docker run -p 8081:8081 tc-golunch-core-service
```

## 📈 Monitoramento

- **Health Check**: `GET /ping`
- **Swagger UI**: `GET /swagger/index.html`
- **Logs**: Estruturados em JSON

## 🔄 CI/CD

O serviço possui pipeline CI/CD separado em duas fases:

### 📋 **Integração Contínua (ci.yaml)**
- **Trigger**: Push/PR para branch master
- **Testes Automatizados**: Execução de testes unitários e BDD
- **Análise de Cobertura**: Meta mínima de 5%
- **Validação**: Verificação de dependências e build

### 🚀 **Deploy Contínuo (cd.yaml)**
- **Trigger**: Push para branch master (após CI)
- **Build Docker**: Geração de imagem para AWS ECR
- **Deploy AWS**: Deploy automático via Helm/Kubernetes
- **Configuração**: Secrets e variáveis de ambiente

## 📝 Documentação da API

A documentação completa da API está disponível via Swagger UI em:
`http://localhost:8081/swagger/index.html`

## 🔗 Integração Serverless (AWS Lambda)

✅ **PRONTO PARA USO**: A autenticação serverless já está totalmente configurada!

### **🛠️ Código Implementado**
O código foi atualizado seguindo o padrão do monolítico `tc-golunch-api`:

1. **ServerlessAuthGateway**: Implementado para comunicação com Lambda
2. **ServerlessAuthMiddleware**: Middleware de autenticação serverless  
3. **main.go**: Atualizado para usar serverless auth em vez de JWT local

### **🔧 Configuração das URLs**

**Apenas configure as URLs serverless** (o resto já está pronto):

```bash
# URLs das funções Lambda (obtidas após deploy do tc-golunch-serverless)
export LAMBDA_AUTH_URL="https://seu-api-gateway-id.execute-api.region.amazonaws.com/auth"
export SERVICE_AUTH_LAMBDA_URL="https://seu-api-gateway-id.execute-api.region.amazonaws.com/service-auth"

# Variáveis existentes (mantidas)
export DATABASE_URL="host=localhost user=golunch_order password=golunch_order123 dbname=golunch_orders port=5433 sslmode=disable TimeZone=America/Sao_Paulo"
export PAYMENT_SERVICE_URL="http://localhost:8082"
export OPERATION_SERVICE_URL="http://localhost:8083"
```

### **📦 Deploy Kubernetes**

✅ **CONFIGURADO**: Os manifestos Kubernetes já estão configurados com as variáveis serverless!

**Apenas ajuste as URLs** no ConfigMap antes do deploy:

```bash
# 1. Edite o ConfigMap com suas URLs reais
vim k8s/core-service-configmap.yaml

# Substitua:
# LAMBDA_AUTH_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/auth"
# SERVICE_AUTH_LAMBDA_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/service-auth"

# 2. Deploy completo
kubectl apply -f k8s/
```

**Estrutura já configurada:**
```yaml
# k8s/core-service-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: core-service-config
data:
  LAMBDA_AUTH_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/auth"
  SERVICE_AUTH_LAMBDA_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/service-auth"
  # ... outras variáveis
```

### **✅ Verificação da Configuração**

Após configurar as variáveis, teste a integração:

```bash
# 1. Inicie o serviço
go run cmd/api/main.go

# 2. Teste autenticação serverless
curl -X POST http://localhost:8081/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 3. Verifique logs para confirmação da integração Lambda
```

### **🔄 Migração Gradual**

A implementação mantém **compatibilidade total** com o código existente:
- ✅ Mesmas interfaces de autenticação
- ✅ Mesmos endpoints e responses
- ✅ Zero breaking changes para clientes
- ✅ Fallback automático se Lambda não disponível
