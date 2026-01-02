package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-order-service/internal/customer/dto"
)

type DataSource interface {
	Create(ctx context.Context, customer dto.CustomerDAO) (dto.CustomerDAO, error)
	FindByCPF(ctx context.Context, cpf string) (dto.CustomerDAO, error)
}
