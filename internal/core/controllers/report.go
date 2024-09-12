package controllers

import (
	"go-template/models"
)

type ReportsController interface {
	ReadReport(id uint) (models.ProductReport, error)
	CreateReport(input models.CreateProductReportInput) (models.ProductReport, error)
}

func (s *ServerImpl) CreateReport(input models.CreateProductReportInput) (models.ProductReport, error) {
	var (
		report models.ProductReport
		err    error
	)

	for i := 0; i < 10; i++ {
		report, err = (*s.Report).WithTrx(s.tx).Create(input)
		if err != nil {
			return report, err
		}
	}

	return report, nil
}

func (s *ServerImpl) ReadReport(id uint) (models.ProductReport, error) {
	report, err := (*s.Report).GetByID(id)
	if err != nil {
		return report, err
	}

	return report, nil
}
