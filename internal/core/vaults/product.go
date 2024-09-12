package vaults

import (
	"go-template/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/vaults_test/$GOFILE -package=vaults_test

type Product interface {
	Create(product *models.Product) error
	GetByID(id uint) (models.Product, error)
}

type ProductImpl struct {
	db *gorm.DB
}

func (p ProductImpl) Create(product *models.Product) error {
	return p.db.Create(product).Error
}

func (p ProductImpl) GetByID(id uint) (product models.Product, err error) {
	err = p.db.First(&product, id).Error

	return
}

func NewProduct(db *gorm.DB) Product {
	return ProductImpl{db: db}
}
