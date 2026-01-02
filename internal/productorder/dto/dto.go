package dto

import (
	"time"

	"github.com/fiap-161/tc-golunch-order-service/internal/productorder/entity"
	coreentity "github.com/fiap-161/tc-golunch-order-service/internal/shared/entity"
	"github.com/google/uuid"
)

type ProductOrderDAO struct {
	coreentity.Entity
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type ProductOrderRequestDTO struct {
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type ProductOrderResponseDTO struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type OrderProductInfo struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Convert entity to DAO
func ToProductOrderDAO(po entity.ProductOrder) ProductOrderDAO {
	return ProductOrderDAO{
		Entity: coreentity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ProductID: po.ProductID,
		OrderID:   po.OrderID,
		Quantity:  po.Quantity,
		UnitPrice: po.UnitPrice,
	}
}

// Convert list of entities to list of DAOs
func ToListProductOrderDAO(list []entity.ProductOrder) []ProductOrderDAO {
	var daos []ProductOrderDAO
	for _, item := range list {
		daos = append(daos, ToProductOrderDAO(item))
	}
	return daos
}

// Convert DAO to entity
func FromProductOrderDAO(dao ProductOrderDAO) entity.ProductOrder {
	return entity.ProductOrder{
		ID:        dao.ID,
		ProductID: dao.ProductID,
		OrderID:   dao.OrderID,
		Quantity:  dao.Quantity,
		UnitPrice: dao.UnitPrice,
	}
}

// Convert list of DAOs to list of entities
func ToListProductOrder(list []ProductOrderDAO) []entity.ProductOrder {
	var result []entity.ProductOrder
	for _, dao := range list {
		result = append(result, FromProductOrderDAO(dao))
	}
	return result
}

// Convert request DTO to entity
func FromRequestDTO(dto ProductOrderRequestDTO) entity.ProductOrder {
	return entity.ProductOrder{
		ProductID: dto.ProductID,
		OrderID:   dto.OrderID,
		Quantity:  dto.Quantity,
		UnitPrice: dto.UnitPrice,
	}
}
