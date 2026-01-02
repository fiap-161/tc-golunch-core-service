package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-core-service/internal/product/dto"
	apperror "github.com/fiap-161/tc-golunch-core-service/internal/shared/errors"
	"gorm.io/gorm"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
	Find(dest any, conds ...any) *gorm.DB
	Delete(value any, conds ...any) *gorm.DB
	Model(value any) *gorm.DB
	Updates(values any) *gorm.DB
}

type GormDataSource struct {
	db DB
}

func New(db DB) *GormDataSource {
	return &GormDataSource{
		db: db,
	}
}

func (r *GormDataSource) Create(_ context.Context, productDAO dto.ProductDAO) (dto.ProductDAO, error) {
	tx := r.db.Create(&productDAO)
	if tx.Error != nil {
		return dto.ProductDAO{}, tx.Error
	}

	return productDAO, nil
}

func (r *GormDataSource) GetAllByCategory(_ context.Context, category string) ([]dto.ProductDAO, error) {
	var products []dto.ProductDAO
	query := r.db
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *GormDataSource) Update(ctx context.Context, id string, updated dto.ProductDAO) (dto.ProductDAO, error) {
	existing, err := r.FindByID(ctx, id)
	if err != nil {
		return dto.ProductDAO{}, err
	}

	updates := map[string]any{}
	if updated.Name != "" {
		updates["name"] = updated.Name
	}
	if updated.Description != "" {
		updates["description"] = updated.Description
	}
	if updated.ImageURL != "" {
		updates["image_url"] = updated.ImageURL
	}
	if updated.Price != 0 {
		updates["price"] = updated.Price
	}
	if updated.PreparingTime != 0 {
		updates["preparing_time"] = updated.PreparingTime
	}
	if updated.Category != "" {
		updates["category"] = updated.Category
	}

	if len(updates) == 0 {
		return existing, nil
	}

	if err := r.db.Model(&dto.ProductDAO{}).Where("id = @id", map[string]any{"id": id}).Updates(updates).Error; err != nil {
		return dto.ProductDAO{}, err
	}

	var updatedProduct dto.ProductDAO
	if err := r.db.Where("id = @id", map[string]any{"id": id}).First(&updatedProduct).Error; err != nil {
		return dto.ProductDAO{}, err
	}

	return updatedProduct, nil
}

func (r *GormDataSource) FindByID(_ context.Context, id string) (dto.ProductDAO, error) {
	var product dto.ProductDAO
	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		if err.Error() == "record not found" {
			return dto.ProductDAO{}, &apperror.NotFoundError{Msg: "Product not found"}
		}
		return dto.ProductDAO{}, err
	}

	return product, nil
}

func (r *GormDataSource) FindByIDs(_ context.Context, ids []string) ([]dto.ProductDAO, error) {
	var products []dto.ProductDAO

	if err := r.db.Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, &apperror.NotFoundError{Msg: "No products found"}
	}

	return products, nil
}

func (r *GormDataSource) Delete(_ context.Context, id string) error {
	var product dto.ProductDAO

	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		if err.Error() == "record not found" {
			return &apperror.NotFoundError{Msg: "Product not found"}
		}
		return err
	}

	if err := r.db.Delete(&product).Error; err != nil {
		return err
	}

	return nil
}
