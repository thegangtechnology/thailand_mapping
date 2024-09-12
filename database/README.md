### Helpers
This folder handles database connection and health check the database.

Import table from `models` in the project

### How to create migration scripts
There is no more auto migrations, sadly. But, here is what you could do.
```zsh
goose -dir=./database/duck/migrations create what_are_you_doing sql
goose -dir=./database/duck/migrations create what_are_you_doing go
```

You can still use gorm in the project.
Inside docker, use `make mig-up`.

### Enums in database
I use `enums.go`, you can rely on that file. It will let you do multiple things to create `USER-DEFINED` type.

### Examples
```go
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

// Define what would you do when you upgrade to this version
func upAddProduct(tx *sql.Tx) error {
	// Step 1 get Gorm tool by connecting to PostgreSQL using sql.Tx
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	// Step 2 do whatever you want with Migrator tool from Gorm
	err = gormDB.Migrator().CreateTable(&models.Product{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}

// Define what would you do when you downgrade from this version
func downAddProduct(tx *sql.Tx) error {
	// Step 1 get Gorm tool by connecting to PostgreSQL using sql.Tx
	gormDB, err := database.GetGormTools(tx)
	if err != nil {
		return tx.Rollback()
	}

	// Step 2 do whatever you want with Migrator tool from Gorm
	err = gormDB.Migrator().DropTable(&models.Product{})
	if err != nil {
		return tx.Rollback()
	}

	return nil
}
```