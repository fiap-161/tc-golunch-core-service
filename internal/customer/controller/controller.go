package controller

import (
	"context"

	"github.com/fiap-161/tc-golunch-order-service/internal/customer/dto"
	"github.com/fiap-161/tc-golunch-order-service/internal/customer/external/datasource"
	"github.com/fiap-161/tc-golunch-order-service/internal/customer/gateway"
	"github.com/fiap-161/tc-golunch-order-service/internal/customer/usecases"
	apperror "github.com/fiap-161/tc-golunch-order-service/internal/shared/errors"
	"github.com/google/uuid"
)

type Controller struct {
	CustomerDataSource datasource.DataSource
	AuthGateway        gateway.AuthGateway
}

func Build(customerDataSource datasource.DataSource, authGateway gateway.AuthGateway) *Controller {
	return &Controller{
		CustomerDataSource: customerDataSource,
		AuthGateway:        authGateway,
	}
}

func (c *Controller) Create(ctx context.Context, customerRequest dto.CustomerRequestDTO) (string, error) {
	customerGateway := gateway.Build(c.CustomerDataSource)
	useCase := usecases.Build(*customerGateway)
	customer := dto.FromCustomerRequestDTO(customerRequest)
	customerId, err := useCase.Create(ctx, customer)

	if err != nil {
		return "", err
	}

	return customerId, nil

}

func (c *Controller) Identify(ctx context.Context, cpf string) (string, error) {

	if cpf == "" {
		return c.createAnonymousToken()
	}

	customerGateway := gateway.Build(c.CustomerDataSource)
	useCase := usecases.Build(*customerGateway)
	customerId, err := useCase.FindByCPF(ctx, cpf)

	if err != nil {
		return "", err
	}

	token, err := c.createToken(customerId, false)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (c *Controller) createAnonymousToken() (string, error) {
	anonymousID := uuid.NewString()

	token, err := c.createToken(anonymousID, true)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *Controller) createToken(id string, isAnonymous bool) (string, error) {
	additionalClaims := map[string]any{
		"is_anonymous": isAnonymous,
	}

	token, err := c.AuthGateway.GenerateToken(id, "customer", additionalClaims)
	if err != nil {
		return "", &apperror.InternalError{Msg: "Error creating token"}
	}

	return token, nil
}
