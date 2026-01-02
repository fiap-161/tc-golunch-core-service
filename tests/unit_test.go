package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockCustomerRepository simula o repositório de clientes
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) FindByCPF(cpf string) (*CustomerDAO, error) {
	args := m.Called(cpf)
	return args.Get(0).(*CustomerDAO), args.Error(1)
}

func (m *MockCustomerRepository) Create(customer *CustomerDAO) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) FindByID(id string) (*CustomerDAO, error) {
	args := m.Called(id)
	return args.Get(0).(*CustomerDAO), args.Error(1)
}

// MockJWTService simula o serviço JWT
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(claims map[string]interface{}) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// CustomerDAO representa a estrutura do cliente no banco
type CustomerDAO struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CPF         string `json:"cpf"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// CustomerResponse representa a resposta da API
type CustomerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	CPF   string `json:"cpf"`
	Token string `json:"token,omitempty"`
}

// OrderRequest representa uma requisição de pedido
type OrderRequest struct {
	CustomerID string `json:"customer_id"`
	Products   []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"products"`
}

// OrderResponse representa a resposta de um pedido
type OrderResponse struct {
	ID            string  `json:"id"`
	CustomerID    string  `json:"customer_id"`
	Status        string  `json:"status"`
	TotalAmount   float64 `json:"total_amount"`
	QRCode        string  `json:"qr_code,omitempty"`
	PreparingTime int     `json:"preparing_time"`
}

// ProductResponse representa a resposta de um produto
type ProductResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

func TestAnonymousCustomerEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock services
	mockJWT := new(MockJWTService)

	// Expected token for anonymous user
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.anonymous"

	// Configure mock expectations
	mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}")).Return(expectedToken, nil)

	// Setup endpoint
	router.GET("/customer/anonymous", func(c *gin.Context) {
		// Simula a geração de token anônimo
		claims := map[string]interface{}{
			"customer_id":  "anonymous_uuid",
			"is_anonymous": true,
		}

		token, err := mockJWT.GenerateToken(claims)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, gin.H{
			"token":        token,
			"customer_id":  "anonymous_uuid",
			"is_anonymous": true,
			"message":      "Anonymous token generated successfully",
		})
	})

	// Execute request
	req, _ := http.NewRequest("GET", "/customer/anonymous", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedToken, response["token"])
	assert.Equal(t, true, response["is_anonymous"])
	assert.Contains(t, response["message"], "Anonymous token generated successfully")

	// Verify mock expectations
	mockJWT.AssertExpectations(t)
}

func TestCustomerRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock services
	mockJWT := new(MockJWTService)
	mockCustomerRepo := new(MockCustomerRepository)

	// Test data
	customerData := CustomerDAO{
		ID:          "customer_123",
		Name:        "João Silva",
		Email:       "joao@test.com",
		CPF:         "12345678901",
		IsAnonymous: false,
	}

	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.customer"

	// Configure mock expectations
	mockCustomerRepo.On("FindByCPF", customerData.CPF).Return((*CustomerDAO)(nil), gorm.ErrRecordNotFound)
	mockCustomerRepo.On("Create", mock.AnythingOfType("*tests.CustomerDAO")).Return(nil)
	mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}")).Return(expectedToken, nil)

	// Setup endpoint
	router.POST("/customer/register", func(c *gin.Context) {
		var customer CustomerDAO
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Check if customer exists
		existingCustomer, err := mockCustomerRepo.FindByCPF(customer.CPF)
		if err != gorm.ErrRecordNotFound && existingCustomer != nil {
			c.JSON(409, gin.H{"error": "Customer already exists"})
			return
		}

		// Create customer
		customer.ID = "customer_123"
		customer.IsAnonymous = false
		err = mockCustomerRepo.Create(&customer)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create customer"})
			return
		}

		// Generate token
		claims := map[string]interface{}{
			"customer_id":  customer.ID,
			"is_anonymous": false,
		}

		token, err := mockJWT.GenerateToken(claims)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(201, CustomerResponse{
			ID:    customer.ID,
			Name:  customer.Name,
			Email: customer.Email,
			CPF:   customer.CPF,
			Token: token,
		})
	})

	// Execute request
	jsonData, _ := json.Marshal(customerData)
	req, _ := http.NewRequest("POST", "/customer/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 201, w.Code)

	var response CustomerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "customer_123", response.ID)
	assert.Equal(t, "João Silva", response.Name)
	assert.Equal(t, expectedToken, response.Token)

	// Verify mock expectations
	mockCustomerRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestCustomerIdentifyByCPF(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock services
	mockJWT := new(MockJWTService)
	mockCustomerRepo := new(MockCustomerRepository)

	// Test data
	customerData := &CustomerDAO{
		ID:          "customer_123",
		Name:        "João Silva",
		Email:       "joao@test.com",
		CPF:         "12345678901",
		IsAnonymous: false,
	}

	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.customer"

	// Configure mock expectations
	mockCustomerRepo.On("FindByCPF", "12345678901").Return(customerData, nil)
	mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}")).Return(expectedToken, nil)

	// Setup endpoint
	router.GET("/customer/identify/:cpf", func(c *gin.Context) {
		cpf := c.Param("cpf")

		customer, err := mockCustomerRepo.FindByCPF(cpf)
		if err != nil {
			c.JSON(404, gin.H{"error": "Customer not found"})
			return
		}

		// Generate token
		claims := map[string]interface{}{
			"customer_id":  customer.ID,
			"is_anonymous": false,
		}

		token, err := mockJWT.GenerateToken(claims)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, CustomerResponse{
			ID:    customer.ID,
			Name:  customer.Name,
			Email: customer.Email,
			CPF:   customer.CPF,
			Token: token,
		})
	})

	// Execute request
	req, _ := http.NewRequest("GET", "/customer/identify/12345678901", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)

	var response CustomerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "customer_123", response.ID)
	assert.Equal(t, "João Silva", response.Name)
	assert.Equal(t, expectedToken, response.Token)

	// Verify mock expectations
	mockCustomerRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestProductListByCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock data
	products := []ProductResponse{
		{
			ID:       "prod_001",
			Name:     "Big GoLunch",
			Price:    25.90,
			Category: "Lanche",
		},
		{
			ID:       "prod_002",
			Name:     "GoLunch Simples",
			Price:    15.90,
			Category: "Lanche",
		},
	}

	// Setup endpoint
	router.GET("/product", func(c *gin.Context) {
		category := c.Query("category")

		var filteredProducts []ProductResponse
		for _, product := range products {
			if category == "" || product.Category == category {
				filteredProducts = append(filteredProducts, product)
			}
		}

		c.JSON(200, filteredProducts)
	})

	// Test 1: List all products
	req, _ := http.NewRequest("GET", "/product", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var allProducts []ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &allProducts)
	assert.NoError(t, err)
	assert.Len(t, allProducts, 2)

	// Test 2: Filter by category
	req, _ = http.NewRequest("GET", "/product?category=Lanche", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var categoryProducts []ProductResponse
	err = json.Unmarshal(w.Body.Bytes(), &categoryProducts)
	assert.NoError(t, err)
	assert.Len(t, categoryProducts, 2)
	assert.Equal(t, "Lanche", categoryProducts[0].Category)
}

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup health endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Execute request
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "pong", response["message"])
}
