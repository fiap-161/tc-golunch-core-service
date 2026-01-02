package presenter

import (
	"github.com/fiap-161/tc-golunch-order-service/internal/product/dto"
	"github.com/fiap-161/tc-golunch-order-service/internal/product/entity"
)

type Presenter struct {
}

func Build() *Presenter {
	return &Presenter{}
}

func (p *Presenter) FromEntityToResponseDTO(product entity.Product) dto.ProductResponseDTO {
	return dto.ProductResponseDTO{
		ID:            product.Id,
		Name:          product.Name,
		Price:         product.Price,
		Description:   product.Description,
		PreparingTime: product.PreparingTime,
		Category:      product.Category,
		ImageURL:      product.ImageURL,
	}
}

func (p *Presenter) FromEntityListToProductListResponseDTO(products []entity.Product) dto.ProductListResponseDTO {
	var productsDTO []dto.ProductResponseDTO
	for _, product := range products {
		dto := p.FromEntityToResponseDTO(product)
		productsDTO = append(productsDTO, dto)
	}

	return dto.ProductListResponseDTO{
		Total: uint(len(productsDTO)),
		List:  productsDTO,
	}
}
