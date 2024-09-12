package migrations

import (
	"database/sql"
	"go-template/database"
	"go-template/models"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddReport, downAddReport)
}

func upAddReport(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	err = gormDB.Migrator().CreateTable(&models.ProductReport{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}

func downAddReport(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	err = gormDB.Migrator().DropTable(&models.ProductReport{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}
