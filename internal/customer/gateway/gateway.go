package gateway

import (
	"context"

	"github.com/fiap-161/tc-golunch-order-service/internal/customer/dto"
	"github.com/fiap-161/tc-golunch-order-service/internal/customer/entity"
	"github.com/fiap-161/tc-golunch-order-service/internal/customer/external/datasource"
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

func (g *Gateway) Create(ctx context.Context, customer entity.Customer) (entity.Customer, error) {
	customerDAO := dto.ToCustomerDAO(customer)

	createdDAO, err := g.datasource.Create(ctx, customerDAO)
	if err != nil {
		return entity.Customer{}, &apperror.InternalError{Msg: err.Error()}
	}

	createdCustomer := dto.FromCustomerDAO(createdDAO)
	return createdCustomer, nil
}

func (g *Gateway) FindByCPF(ctx context.Context, cpf string) (entity.Customer, error) {
	customerDAO, err := g.datasource.FindByCPF(ctx, cpf)
	if err != nil {
		return entity.Customer{}, &apperror.InternalError{Msg: err.Error()}
	}

	customer := dto.FromCustomerDAO(customerDAO)
	return customer, nil
}
