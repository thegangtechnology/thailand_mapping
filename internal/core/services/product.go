package services

import (
	"go-template/internal/core/vaults"
	"go-template/models"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/services_test/$GOFILE -package=services_test

type Product interface {
	Create(input models.CreateProductInput) (models.Product, error)
	GetByID(id uint) (models.Product, error)
}

type ProductImpl struct {
	product vaults.Product
}

func (r ProductImpl) Create(input models.CreateProductInput) (models.Product, error) {
	product := models.Product{
		ProductDTO: input.ProductDTO,
	}

	if err := r.product.Create(&product); err != nil {
		return product, err
	}

	return product, nil
}

func (r ProductImpl) GetByID(id uint) (models.Product, error) {
	return r.product.GetByID(id)
}

func NewProduct(rr vaults.Product) Product {
	return ProductImpl{product: rr}
}
