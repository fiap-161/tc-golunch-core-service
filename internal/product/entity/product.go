package entity

import (
	"github.com/fiap-161/tc-golunch-core-service/internal/product/entity/enum"
	apperror "github.com/fiap-161/tc-golunch-core-service/internal/shared/errors"
)

type Product struct {
	Id            string
	Name          string
	Price         float64
	Description   string
	PreparingTime uint
	Category      enum.Category
	ImageURL      string
}

func (p Product) Build() Product {
	return Product{
		Id:            p.Id,
		Name:          p.Name,
		Price:         p.Price,
		Description:   p.Description,
		PreparingTime: p.PreparingTime,
		Category:      p.Category,
		ImageURL:      p.ImageURL,
	}
}

func (p Product) Validate() error {
	if p.Name == "" {
		return &apperror.ValidationError{Msg: "Name is required"}
	}
	if p.Price < 0 {
		return &apperror.ValidationError{Msg: "Price must be positive"}
	}
	if p.Category == "" {
		return &apperror.ValidationError{Msg: "Category is required"}
	}
	return nil
}
