package controller

import (
	"context"
	"mime/multipart"

	"github.com/fiap-161/tc-golunch-core-service/internal/product/dto"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/entity"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/entity/enum"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/external/datasource"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/gateway"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/presenter"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/usecases"
)

// controller (DTOS)

// INFO: controllers Criam gateways e requisitam usecases
type Controller struct {
	productDatasource datasource.DataSource
}

func Build(productDataSource datasource.DataSource) *Controller {
	return &Controller{
		productDatasource: productDataSource,
	}
}

func (c *Controller) Create(ctx context.Context, productDTO dto.ProductRequestDTO) (dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	var product entity.Product
	product = dto.FromRequestDTO(productDTO)
	createdProduct, createErr := useCase.CreateProduct(ctx, product)
	if createErr != nil {
		return dto.ProductResponseDTO{}, createErr
	}

	return presenter.FromEntityToResponseDTO(createdProduct), nil
}

func (c *Controller) ListCategories(ctx context.Context) []enum.Category {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)
	return useCase.ListCategories(ctx)
}

func (c *Controller) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)
	return useCase.UploadImage(ctx, fileHeader)
}

func (c *Controller) GetAllByCategory(ctx context.Context, category string) (dto.ProductListResponseDTO, error) {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	result, err := useCase.GetAllByCategory(ctx, category)

	if err != nil {
		return dto.ProductListResponseDTO{}, err
	}

	return presenter.FromEntityListToProductListResponseDTO(result), nil
}

func (c *Controller) Update(ctx context.Context, productId string, productDTO dto.ProductRequestUpdateDTO) (dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	product := dto.FromUpdateDTO(productDTO)
	result, err := useCase.Update(ctx, productId, product)
	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	return presenter.FromEntityToResponseDTO(result), nil
}

func (c *Controller) FindByID(ctx context.Context, productId string) (dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	result, err := useCase.FindByID(ctx, productId)
	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	return presenter.FromEntityToResponseDTO(result), nil
}

func (c *Controller) Delete(ctx context.Context, productId string) error {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)

	err := useCase.Delete(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) FindByIDs(ctx context.Context, productIdList []string) ([]dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.productDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	result, err := useCase.FindByIDs(ctx, productIdList)
	if err != nil {
		return []dto.ProductResponseDTO{}, err
	}

	return presenter.FromEntityListToProductListResponseDTO(result).List, nil
}
