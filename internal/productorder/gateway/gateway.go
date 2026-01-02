package gateway

import (
	"context"

	"github.com/fiap-161/tc-golunch-order-service/internal/productorder/dto"
	"github.com/fiap-161/tc-golunch-order-service/internal/productorder/entity"
	"github.com/fiap-161/tc-golunch-order-service/internal/productorder/external/datasource"
	apperror "github.com/fiap-161/tc-golunch-order-service/internal/shared/errors"
)

type Gateway struct {
	datasource datasource.DataSource
}

func Build(datasource datasource.DataSource) *Gateway {
	return &Gateway{
		datasource: datasource,
	}
}

func (g *Gateway) CreateBulk(c context.Context, productOrders []entity.ProductOrder) (int, error) {
	var listProductOrderDAO = dto.ToListProductOrderDAO(productOrders)
	length, err := g.datasource.CreateBulk(c, listProductOrderDAO)

	if err != nil {
		return 0, &apperror.InternalError{Msg: err.Error()}
	}

	return length, nil
}

func (g *Gateway) FindByOrderID(c context.Context, orderId string) ([]entity.ProductOrder, error) {
	listProductOrderFoundDAO, err := g.datasource.FindByOrderID(c, orderId)
	productOrder := dto.ToListProductOrder(listProductOrderFoundDAO)

	if err != nil {
		return []entity.ProductOrder{}, &apperror.InternalError{Msg: err.Error()}
	}

	return productOrder, nil
}
