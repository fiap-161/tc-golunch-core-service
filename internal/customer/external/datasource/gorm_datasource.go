package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-core-service/internal/customer/dto"
	"gorm.io/gorm"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
}

type GormDataSource struct {
	db DB
}

func New(db DB) *GormDataSource {
	return &GormDataSource{
		db: db,
	}
}

func (r *GormDataSource) Create(_ context.Context, customer dto.CustomerDAO) (dto.CustomerDAO, error) {
	tx := r.db.Create(&customer)
	if tx.Error != nil {
		return dto.CustomerDAO{}, tx.Error
	}

	return customer, nil
}

func (r *GormDataSource) FindByCPF(_ context.Context, cpf string) (dto.CustomerDAO, error) {
	var customer dto.CustomerDAO

	tx := r.db.Where("cpf = ?", cpf).First(&customer)
	if tx.Error != nil {
		return dto.CustomerDAO{}, tx.Error
	}

	return customer, nil
}
