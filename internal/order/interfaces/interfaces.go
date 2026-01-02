package interfaces

import (
	"context"

	productentity "github.com/fiap-161/tc-golunch-core-service/internal/product/entity"
	productorderentity "github.com/fiap-161/tc-golunch-core-service/internal/productorder/entity"
)

type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]productentity.Product, error)
}

type ProductOrderService interface {
	CreateBulk(ctx context.Context, productOrders []productorderentity.ProductOrder) (int, error)
}
