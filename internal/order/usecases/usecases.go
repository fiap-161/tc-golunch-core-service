package usecases

import (
	"context"

	"github.com/fiap-161/tc-golunch-order-service/internal/order/dto"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/entity"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/gateway"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/interfaces"
	productentity "github.com/fiap-161/tc-golunch-order-service/internal/product/entity"
	productorderentity "github.com/fiap-161/tc-golunch-order-service/internal/productorder/entity"
	apperror "github.com/fiap-161/tc-golunch-order-service/internal/shared/errors"
)

type UseCases struct {
	orderGateway        *gateway.Gateway
	productService      interfaces.ProductService
	productOrderService interfaces.ProductOrderService
}

func Build(
	orderGateway *gateway.Gateway,
	productService interfaces.ProductService,
	productOrderService interfaces.ProductOrderService,
) *UseCases {
	return &UseCases{
		orderGateway:        orderGateway,
		productService:      productService,
		productOrderService: productOrderService,
	}
}

func (u *UseCases) CreateCompleteOrder(ctx context.Context, orderDTO dto.CreateOrderDTO) (entity.Order, error) {
	var productIds []string
	for _, item := range orderDTO.Products {
		productIds = append(productIds, item.ProductID)
	}

	products, findErr := u.productService.FindByIDs(ctx, productIds)
	if findErr != nil {
		return entity.Order{}, findErr
	}
	if len(products) != len(orderDTO.Products) {
		return entity.Order{}, &apperror.NotFoundError{
			Msg: "some products not found",
		}
	}

	populatedOrder := generateOrderByProducts(orderDTO, products)
	createdOrder, createErr := u.orderGateway.Create(ctx, populatedOrder.Build())
	if createErr != nil {
		return entity.Order{}, createErr
	}

	productOrders, _ := generateProductOrderFromOrderAndProducts(createdOrder.ID, orderDTO.Products, products)
	_, createBulkErr := u.productOrderService.CreateBulk(ctx, productOrders)
	if createBulkErr != nil {
		return entity.Order{}, createBulkErr
	}

	return createdOrder, nil
}

func generateOrderByProducts(orderDTO dto.CreateOrderDTO, products []productentity.Product) entity.Order {
	orderProductInfo := make([]entity.OrderProductInfo, len(orderDTO.Products))
	for i, product := range orderDTO.Products {
		orderProductInfo[i] = entity.OrderProductInfo{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		}
	}

	return entity.Order{}.FromDTO(orderDTO.CustomerID, orderProductInfo, products)
}

func generateProductOrderFromOrderAndProducts(
	orderID string,
	orderProductInfo []dto.OrderProductInfo,
	products []productentity.Product,
) ([]productorderentity.ProductOrder, error) {
	var res []productorderentity.ProductOrder

	for _, product := range products {
		for _, item := range orderProductInfo {
			if product.Id == item.ProductID {
				res = append(res, productorderentity.ProductOrder{
					OrderID:   orderID,
					ProductID: product.Id,
					Quantity:  item.Quantity,
					UnitPrice: product.Price,
				})
			}
		}
	}

	return res, nil
}

func (u *UseCases) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.orderGateway.Create(ctx, order)
}

func (u *UseCases) GetAllOrById(ctx context.Context, id string) ([]entity.Order, error) {
	if id != "" {
		order, err := u.orderGateway.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return []entity.Order{order}, nil
	}
	return u.orderGateway.GetAll(ctx)
}

func (u *UseCases) GetPanel(ctx context.Context) ([]entity.Order, error) {
	return u.orderGateway.GetPanel(ctx)
}

func (u *UseCases) FindByID(ctx context.Context, id string) (entity.Order, error) {
	return u.orderGateway.FindByID(ctx, id)
}

func (u *UseCases) Update(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.orderGateway.Update(ctx, order)
}
