package duck

//nolint:revive // goose use this
import (
	"embed"
	"flag"
	"go-template/database"
	_ "go-template/database/duck/migrations"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

//go:embed migrations/*.go
var embedMigrations embed.FS

func UseGoose(args []string) {
	dir := flag.String("dir", "migrations", "directory with migration files")

	var arguments []string

	if len(args) == 0 {
		return
	}

	goose.SetBaseFS(embedMigrations)

	command := args[0]
	arguments = append(arguments, args[1:]...)

	defer func() {
		if r := recover(); r != nil {
			_ = database.Close()
		}
	}()

	db, err := database.Get().DB()
	if err != nil {
		log.WithError(err).Panic(args)
	}

	if err = goose.SetDialect("postgres"); err != nil {
		log.WithError(err).Panic(args)
	}

	if err = goose.Run(command, db, *dir, arguments...); err != nil {
		log.Panicf("goose run: %v", err)
	}
}
