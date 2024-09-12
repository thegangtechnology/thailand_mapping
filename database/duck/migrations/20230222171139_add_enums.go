package migrations

import (
	"database/sql"
	"go-template/database"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddEnums, downAddEnums)
}

func upAddEnums(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	if !checkEnumExist(gormDB, "transaction_action_enum") {
		createEnum(
			gormDB,
			"transaction_action_enum",
			[]string{
				"update",
				"create",
				"delete",
				"revert",
				"submit",
				"autofill",
				"manual",
			},
		)
	}

	if !checkEnumExist(gormDB, "report_status_enum") {
		createEnum(
			gormDB,
			"report_status_enum",
			[]string{
				"pending",
				"draft",
				"completed",
				"deleted",
			},
		)
	}

	if !checkEnumExist(gormDB, "well_status_enum") {
		createEnum(
			gormDB,
			"well_status_enum",
			[]string{
				"closed",
				"opened",
			},
		)
	}

	return nil
}

func downAddEnums(tx *sql.Tx) error {
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	removeEnum(
		gormDB,
		"transaction_action_enum",
	)
	removeEnum(
		gormDB,
		"report_status_enum",
	)
	removeEnum(
		gormDB,
		"well_status_enum",
	)

	return nil
}
