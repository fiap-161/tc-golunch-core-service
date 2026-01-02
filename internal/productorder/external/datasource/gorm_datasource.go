package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-order-service/internal/productorder/dto"
	"gorm.io/gorm"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	Find(dest any, conds ...any) *gorm.DB
}

type GormDataSource struct {
	db DB
}

func New(db DB) *GormDataSource {
	return &GormDataSource{
		db: db,
	}
}

func (r *GormDataSource) CreateBulk(_ context.Context, orders []dto.ProductOrderDAO) (int, error) {
	tx := r.db.Create(&orders)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return len(orders), nil
}

func (r *GormDataSource) FindByOrderID(_ context.Context, orderID string) ([]dto.ProductOrderDAO, error) {
	var orders []dto.ProductOrderDAO

	tx := r.db.Where("order_id = ?", orderID).Find(&orders)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return orders, nil
}
