package dto

import (
	"strings"
	"time"

	"github.com/fiap-161/tc-golunch-order-service/internal/product/entity"
	"github.com/fiap-161/tc-golunch-order-service/internal/product/entity/enum"
	coreentity "github.com/fiap-161/tc-golunch-order-service/internal/shared/entity"
	"github.com/google/uuid"
)

type ProductRequestDTO struct {
	Name          string        `json:"name" binding:"required"`
	Price         float64       `json:"price" binding:"required"`
	Description   string        `json:"description" binding:"required"`
	PreparingTime uint          `json:"preparing_time" binding:"required"`
	Category      enum.Category `json:"category" binding:"required"`
	ImageURL      string        `json:"image_url" binding:"required,url"`
}

type ProductResponseDTO struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Price         float64       `json:"price"`
	Description   string        `json:"description"`
	PreparingTime uint          `json:"preparing_time"`
	Category      enum.Category `json:"category"`
	ImageURL      string        `json:"image_url"`
}

type ProductListResponseDTO struct {
	Total uint                 `json:"total"`
	List  []ProductResponseDTO `json:"list"`
}

type ProductRequestUpdateDTO struct {
	Name          string        `json:"name"`
	Price         float64       `json:"price"`
	Description   string        `json:"description"`
	PreparingTime uint          `json:"preparing_time"`
	Category      enum.Category `json:"category"`
	ImageURL      string        `json:"image_url"`
}

type ImageURLDTO struct {
	ImageURL string `json:"url"`
}

type ProductDAO struct {
	coreentity.Entity
	Name          string        `json:"name"`
	Price         float64       `json:"price" gorm:"type:decimal(10,2)"`
	Description   string        `json:"description" gorm:"type:text"`
	PreparingTime uint          `json:"preparing_time" gorm:"type:integer"`
	Category      enum.Category `json:"category"`
	ImageURL      string        `json:"image_url" gorm:"type:varchar(255)"`
}

// Convert entity entity to DAO
func ToProductDAO(p entity.Product) ProductDAO {
	return ProductDAO{
		Entity: coreentity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:          p.Name,
		Price:         p.Price,
		Description:   p.Description,
		PreparingTime: p.PreparingTime,
		Category:      p.Category,
		ImageURL:      p.ImageURL,
	}
}

// Convert DAO to entity entity
func FromProductDAO(dao ProductDAO) entity.Product {
	category := strings.ToUpper(string(dao.Category))
	return entity.Product{
		Id:            dao.ID,
		Name:          dao.Name,
		Price:         dao.Price,
		Description:   dao.Description,
		PreparingTime: dao.PreparingTime,
		Category:      enum.Category(category),
		ImageURL:      dao.ImageURL,
	}
}

// Convert request DTO to entity entity
func FromRequestDTO(dto ProductRequestDTO) entity.Product {
	category := strings.ToUpper(string(dto.Category))
	return entity.Product{
		Name:          dto.Name,
		Price:         dto.Price,
		Description:   dto.Description,
		PreparingTime: dto.PreparingTime,
		Category:      enum.Category(category),
		ImageURL:      dto.ImageURL,
	}
}

// Convert update request DTO to entity entity
func FromUpdateDTO(dto ProductRequestUpdateDTO) entity.Product {
	category := strings.ToUpper(string(dto.Category))
	return entity.Product{
		Name:          dto.Name,
		Price:         dto.Price,
		Description:   dto.Description,
		PreparingTime: dto.PreparingTime,
		Category:      enum.Category(category),
		ImageURL:      dto.ImageURL,
	}
}

func EntityListFromDAOList(daoList []ProductDAO) []entity.Product {
	products := make([]entity.Product, 0, len(daoList))
	for _, dao := range daoList {
		products = append(products, FromProductDAO(dao))
	}
	return products
}
