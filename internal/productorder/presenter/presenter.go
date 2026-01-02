package presenter

import (
	"github.com/fiap-161/tc-golunch-core-service/internal/productorder/dto"
	"github.com/fiap-161/tc-golunch-core-service/internal/productorder/entity"
)

type Presenter struct {
}

func Build() *Presenter {
	return &Presenter{}
}

func (p *Presenter) FromEntityToResponseDTO(po entity.ProductOrder) dto.ProductOrderResponseDTO {
	return dto.ProductOrderResponseDTO{
		ID:        po.ID,
		ProductID: po.ProductID,
		OrderID:   po.OrderID,
		Quantity:  po.Quantity,
		UnitPrice: po.UnitPrice,
	}
}

func (p *Presenter) FromEntityListToResponseDTOList(list []entity.ProductOrder) []dto.ProductOrderResponseDTO {
	var listProductOrderResponseDTO []dto.ProductOrderResponseDTO
	for _, item := range list {
		dto := p.FromEntityToResponseDTO(item)
		listProductOrderResponseDTO = append(listProductOrderResponseDTO, dto)
	}
	return listProductOrderResponseDTO
}
