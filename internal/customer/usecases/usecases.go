package usecases

import (
	"context"

	"github.com/fiap-161/tc-golunch-order-service/internal/customer/entity"
	"github.com/fiap-161/tc-golunch-order-service/internal/customer/gateway"
	apperror "github.com/fiap-161/tc-golunch-order-service/internal/shared/errors"
)

type CustomerUseCases struct {
	CustomerGateway gateway.Gateway
}

func Build(gateway gateway.Gateway) *CustomerUseCases {
	return &CustomerUseCases{
		CustomerGateway: gateway,
	}
}

func (u *CustomerUseCases) Create(ctx context.Context, customer entity.Customer) (string, error) {
	exists, _ := u.CustomerGateway.FindByCPF(ctx, customer.CPF)
	if exists.CPF != "" {
		return "", &apperror.ValidationError{Msg: "Customer already registered"}
	}

	customerWithID := customer.Build()
	saved, err := u.CustomerGateway.Create(ctx, customerWithID)
	if err != nil {
		return "", err
	}

	return saved.Id, nil
}

func (u *CustomerUseCases) FindByCPF(ctx context.Context, cpf string) (string, error) {
	customer, err := u.CustomerGateway.FindByCPF(ctx, cpf)
	if err != nil || customer.Id == "" {
		return "", &apperror.NotFoundError{Msg: "Customer not found"}
	}

	return customer.Id, nil
}
