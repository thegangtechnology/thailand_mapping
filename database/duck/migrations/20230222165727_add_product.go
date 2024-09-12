package migrations

import (
	"database/sql"
	"go-template/database"
	"go-template/models"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddProduct, downAddProduct)
}

func upAddProduct(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	err = gormDB.Migrator().CreateTable(&models.Product{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}

func downAddProduct(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	err = gormDB.Migrator().DropTable(&models.Product{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}
