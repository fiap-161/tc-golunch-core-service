package usecases

import (
	"context"
	"fmt"

	"github.com/fiap-161/tc-golunch-order-service/internal/productorder/entity"
	"github.com/fiap-161/tc-golunch-order-service/internal/productorder/gateway"
	apperror "github.com/fiap-161/tc-golunch-order-service/internal/shared/errors"
)

type UseCases struct {
	productOrderGateway gateway.Gateway
}

func Build(productOrderGateway gateway.Gateway) *UseCases {
	return &UseCases{
		productOrderGateway: productOrderGateway,
	}
}

func (u *UseCases) CreateBulk(ctx context.Context, productOrders []entity.ProductOrder) (int, error) {
	for i, po := range productOrders {
		if po.ProductID == "" {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: productID should not be empty", i)}
		}
		if po.OrderID == "" {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: orderID should not be empty", i)}
		}
		if po.Quantity <= 0 {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: quantity has to be more than zero", i)}
		}
		if po.UnitPrice < 0 {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: unitPrice cannot be negative", i)}
		}
	}

	length, err := u.productOrderGateway.CreateBulk(ctx, productOrders)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (u *UseCases) FindByOrderID(ctx context.Context, orderID string) ([]entity.ProductOrder, error) {
	productOrderFound, err := u.productOrderGateway.FindByOrderID(ctx, orderID)
	if err != nil {
		return []entity.ProductOrder{}, err
	}

	return productOrderFound, nil
}
