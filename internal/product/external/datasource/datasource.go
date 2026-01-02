package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-core-service/internal/product/dto"
)

// TODO: Como é mundo externo, precisa ser uma DTO
type DataSource interface {
	Create(context.Context, dto.ProductDAO) (dto.ProductDAO, error)
	GetAllByCategory(context.Context, string) ([]dto.ProductDAO, error)
	Update(context.Context, string, dto.ProductDAO) (dto.ProductDAO, error)
	FindByID(context.Context, string) (dto.ProductDAO, error)
	FindByIDs(context.Context, []string) ([]dto.ProductDAO, error)
	Delete(context.Context, string) error
}
