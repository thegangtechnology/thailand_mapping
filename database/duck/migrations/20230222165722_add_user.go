package migrations

import (
	"database/sql"
	"go-template/database"
	"go-template/models"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddUser, downAddUser)
}

func upAddUser(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	err = gormDB.Migrator().CreateTable(&models.User{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}

func downAddUser(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	err = gormDB.Migrator().DropTable(&models.User{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}
