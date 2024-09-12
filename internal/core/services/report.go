package services

import (
	"go-template/internal/core/vaults"
	"go-template/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/services_test/$GOFILE -package=services_test

type Report interface {
	Create(input models.CreateProductReportInput) (models.ProductReport, error)
	GetByID(id uint) (models.ProductReport, error)

	WithTrx(db *gorm.DB) Report
}

type ReportImpl struct {
	report vaults.Report
}

func (r ReportImpl) WithTrx(db *gorm.DB) Report {
	r.report = r.report.WithTrx(db)

	return r
}

func (r ReportImpl) Create(input models.CreateProductReportInput) (models.ProductReport, error) {
	report := models.ProductReport{
		ProductReportDTO: input.ProductReportDTO,
	}

	if err := r.report.Create(&report); err != nil {
		return report, err
	}

	return report, nil
}

func (r ReportImpl) GetByID(id uint) (models.ProductReport, error) {
	return r.report.GetByID(id)
}

func NewReport(rr vaults.Report) Report {
	return ReportImpl{report: rr}
}
