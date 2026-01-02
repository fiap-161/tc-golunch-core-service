package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/fiap-161/tc-golunch-core-service/database"
	_ "github.com/fiap-161/tc-golunch-core-service/docs"

	// "github.com/fiap-161/tc-golunch-core-service/internal/http/middleware"
	admincontroller "github.com/fiap-161/tc-golunch-core-service/internal/admin/controller"
	adminmodel "github.com/fiap-161/tc-golunch-core-service/internal/admin/dto"
	admindatasource "github.com/fiap-161/tc-golunch-core-service/internal/admin/external/datasource"
	adminhandler "github.com/fiap-161/tc-golunch-core-service/internal/admin/handler"
	customercontroller "github.com/fiap-161/tc-golunch-core-service/internal/customer/controller"
	customermodel "github.com/fiap-161/tc-golunch-core-service/internal/customer/dto"
	customerdatasource "github.com/fiap-161/tc-golunch-core-service/internal/customer/external/datasource"
	customergateway "github.com/fiap-161/tc-golunch-core-service/internal/customer/gateway"
	customerhandler "github.com/fiap-161/tc-golunch-core-service/internal/customer/handler"
	ordercontroller "github.com/fiap-161/tc-golunch-core-service/internal/order/controller"
	ordermodel "github.com/fiap-161/tc-golunch-core-service/internal/order/dto"
	orderdatasource "github.com/fiap-161/tc-golunch-core-service/internal/order/external/datasource"
	ordergateway "github.com/fiap-161/tc-golunch-core-service/internal/order/gateway"
	orderhandler "github.com/fiap-161/tc-golunch-core-service/internal/order/handler"
	productcontroller "github.com/fiap-161/tc-golunch-core-service/internal/product/controller"
	productmodel "github.com/fiap-161/tc-golunch-core-service/internal/product/dto"
	productdatasource "github.com/fiap-161/tc-golunch-core-service/internal/product/external/datasource"
	productgateway "github.com/fiap-161/tc-golunch-core-service/internal/product/gateway"
	producthandler "github.com/fiap-161/tc-golunch-core-service/internal/product/handler"
	productusecases "github.com/fiap-161/tc-golunch-core-service/internal/product/usecases"
	productordermodel "github.com/fiap-161/tc-golunch-core-service/internal/productorder/dto"
	productorderdatasource "github.com/fiap-161/tc-golunch-core-service/internal/productorder/external/datasource"
	productordergateway "github.com/fiap-161/tc-golunch-core-service/internal/productorder/gateway"
	productorderusecases "github.com/fiap-161/tc-golunch-core-service/internal/productorder/usecases"
	sharedgateway "github.com/fiap-161/tc-golunch-core-service/internal/shared/gateway"
)

// @title           GoLunch Core Service API
// @version         1.0
// @description     API para gerenciamento de pedidos da lanchonete
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host            localhost:8081
// @BasePath        /
func main() {
	r := gin.Default()
	loadYAML()

	db := database.NewPostgresDatabase().GetDb()

	if err := db.AutoMigrate(
		&customermodel.CustomerDAO{},
		&productmodel.ProductDAO{},
		&ordermodel.OrderDAO{},
		&productordermodel.ProductOrderDAO{},
		&adminmodel.AdminDAO{},
	); err != nil {
		log.Fatalf("Erro ao migrar o banco: %v", err)
	}

	// Serverless Auth Gateway (following tc-golunch-api monolith pattern)
	serverlessAuth := sharedgateway.NewServerlessAuthGateway(
		os.Getenv("LAMBDA_AUTH_URL"),
		os.Getenv("SERVICE_AUTH_LAMBDA_URL"),
	)

	// Customer (using serverless auth)
	customerDatasource := customerdatasource.New(db)
	customerController := customercontroller.Build(customerDatasource, serverlessAuth)
	customerHandler := customerhandler.New(customerController)

	// Admin (using same serverless auth as customers - following monolith pattern)
	adminDatasource := admindatasource.New(db)
	adminController := admincontroller.Build(adminDatasource, serverlessAuth)
	adminHandler := adminhandler.New(adminController)

	// Product
	productDataSource := productdatasource.New(db)
	productController := productcontroller.Build(productDataSource)
	productHandler := producthandler.New(productController)

	// Product Order
	productOrderDataSource := productorderdatasource.New(db)
	productOrderGateway := productordergateway.Build(productOrderDataSource)
	productOrderUseCase := productorderusecases.Build(*productOrderGateway)

	// Common Gateways
	productGateway := productgateway.Build(productDataSource)
	productUseCase := productusecases.Build(*productGateway)

	// Order Data Source and Gateway
	orderDataSource := orderdatasource.New(db)
	orderGateway := ordergateway.Build(orderDataSource)

	// Order Controller and Handler
	orderController := ordercontroller.Build(orderGateway, productUseCase, productOrderUseCase)
	orderHandler := orderhandler.New(orderController)

	// Default Routes
	r.GET("/ping", ping)
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	// Public Routes
	r.GET("/customer/identify/:cpf", customerHandler.Identify)
	r.GET("/customer/anonymous", customerHandler.Anonymous)
	r.POST("/customer/register", customerHandler.Create)

	// Admin Routes (centralized authentication)
	r.POST("/admin/register", adminHandler.Register)
	r.POST("/admin/login", adminHandler.Login)
	// r.GET("/admin/validate", adminHandler.ValidateToken) // TODO: Implement validation

	// Product Routes (read-only for core service)
	r.GET("/product/categories", productHandler.ListCategories)
	r.GET("/product", productHandler.GetAllByCategory)

	// Admin Product Routes (protected by JWT token from admin login)
	adminRoutes := r.Group("/admin")
	// adminRoutes.Use(middleware.AdminAuthMiddleware()) // TODO: Implement JWT validation
	adminRoutes.POST("/product", productHandler.Create)
	adminRoutes.PUT("/product/:id", productHandler.Update)
	adminRoutes.DELETE("/product/:id", productHandler.Delete)
	adminRoutes.POST("/product/upload", productHandler.UploadImage)
	adminRoutes.GET("/product", productHandler.GetAllByCategory) // Lista todos os produtos por categoria

	// Order Routes
	r.POST("/order", orderHandler.Create)
	r.GET("/order", orderHandler.GetAll)
	r.PUT("/order/:id", orderHandler.Update)
	r.GET("/order/panel", orderHandler.GetPanel)

	// Webhook Routes for inter-service communication
	webhookHandler := orderhandler.NewWebhookHandler(orderController)
	r.POST("/webhook/payment", webhookHandler.PaymentWebhook)

	r.Run(":8081")
}

func loadYAML() {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf/environment")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading yaml config: %v", err)
	}
}

// Ping godoc
// @Summary      Answers with "pong"
// @Description  Health Check
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200 {object}  PongResponse
// @Router       /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type PongResponse struct {
	Message string `json:"message"`
}
