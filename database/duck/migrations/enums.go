package migrations

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func enumValues(values []string) string {
	trimmedValues := make([]string, len(values))

	for i, v := range values {
		trimmedValues[i] = strings.TrimSpace(v)
	}

	enumValuesString := "'" + strings.Join(trimmedValues, "','") + "'"

	return enumValuesString
}

func checkEnumExist(db *gorm.DB, name string) bool {
	var count int64

	db.Raw("SELECT COUNT(*) FROM pg_type WHERE typname = ?", name).Scan(&count)

	return count > 0
}

func createEnum(db *gorm.DB, name string, enums []string) {
	tx := db.Set("gorm:table_options", "").
		Exec(fmt.Sprintf(`CREATE TYPE %s AS ENUM(%s)`, name, enumValues(enums)))

	if tx.Error != nil {
		tx.Rollback()
		log.Fatalln(fmt.Sprintf("error creating %s:", name), tx.Error)
	}
}

func removeEnum(db *gorm.DB, name string) {
	tx := db.Set("gorm:table_options", "").
		Exec(fmt.Sprintf(`DROP TYPE IF EXISTS %s CASCADE`, name))

	if tx.Error != nil {
		tx.Rollback()
		log.Fatalln(fmt.Sprintf("error creating %s:", name), tx.Error)
	}
}
