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
   export DATABASE_URL="postgres://user:password@localhost:5432/golunch_orders?sslmode=disable"
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
