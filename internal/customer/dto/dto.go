package dto

import (
	"time"

	"github.com/fiap-161/tc-golunch-order-service/internal/customer/entity"
	gormEntity "github.com/fiap-161/tc-golunch-order-service/internal/shared/entity"
	"github.com/google/uuid"
)

type CustomerDAO struct {
	gormEntity.Entity
	Name        string `json:"name"`
	Email       string `json:"email"`
	CPF         string `gorm:"uniqueIndex" json:"cpf"`
	IsAnonymous bool   `json:"is_anonymous" gorm:"default:false"`
}

func ToCustomerDAO(c entity.Customer) CustomerDAO {
	return CustomerDAO{
		Entity: gormEntity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        c.Name,
		Email:       c.Email,
		CPF:         c.CPF,
		IsAnonymous: c.IsAnonymous,
	}
}

func FromCustomerDAO(d CustomerDAO) entity.Customer {
	return entity.Customer{
		Id:          d.ID,
		Name:        d.Name,
		Email:       d.Email,
		CPF:         d.CPF,
		IsAnonymous: d.IsAnonymous,
	}
}

func FromCustomerRequestDTO(dto CustomerRequestDTO) entity.Customer {
	return entity.Customer{
		Name:        dto.Name,
		Email:       dto.Email,
		CPF:         dto.CPF,
		IsAnonymous: false,
	}
}

type CustomerRequestDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	CPF   string `json:"cpf"`
}

type TokenDTO struct {
	TokenString string `json:"token"`
}
