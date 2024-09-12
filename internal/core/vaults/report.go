package vaults

import (
	"go-template/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/vaults_test/$GOFILE -package=vaults_test

type Report interface {
	Create(report *models.ProductReport) error
	GetByID(id uint) (models.ProductReport, error)

	WithTrx(db *gorm.DB) Report
}

type ReportImpl struct {
	db *gorm.DB
}

func (r ReportImpl) WithTrx(db *gorm.DB) Report {
	r.db = db

	return r
}

func (r ReportImpl) Create(report *models.ProductReport) error {
	return r.db.Create(report).Error
}

func (r ReportImpl) GetByID(id uint) (report models.ProductReport, err error) {
	err = r.db.First(&report, id).Error

	return
}

func NewReport(db *gorm.DB) Report {
	return ReportImpl{db: db}
}
