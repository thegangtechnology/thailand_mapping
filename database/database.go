package database

import (
	"database/sql"
	"fmt"
	"go-template/config"
	"go-template/utils/slices"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm/logger"

	log "github.com/sirupsen/logrus"

	"github.com/go-logfmt/logfmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//nolint:gochecknoglobals // serve as the database connection
var instance *gorm.DB

const MaxIdleConns = 10
const MaxOpenConns = 100

type ConnectParams struct {
	ConfigOverride *postgres.Config
}

func ParseDSN(dsn string) map[string]string {
	var values = make(map[string]string)

	d := logfmt.NewDecoder(strings.NewReader(dsn))
	for d.ScanRecord() {
		for d.ScanKeyval() {
			values[string(d.Key())] = string(d.Value())
		}
	}

	if d.Err() != nil {
		panic(d.Err())
	}

	if _, ok := values["host"]; !ok {
		panic(config.ErrHostValueMissing)
	}

	if _, ok := values["user"]; !ok {
		panic(config.ErrUserValueMissing)
	}

	if _, ok := values["password"]; !ok {
		panic(config.ErrPasswordValueMissing)
	}

	if _, ok := values["dbname"]; !ok {
		panic(config.ErrDBNameValueMissing)
	}

	if _, ok := values["port"]; !ok {
		panic(config.ErrPortValueMissing)
	}

	return values
}

func RunDatabaseCommand(command string) error {
	pgConfig := postgres.Config{
		DSN:                  os.Getenv("DB_DSN"),
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	db = db.Exec(command)
	if db.Error != nil {
		return db.Error
	}

	conn, err := db.DB()
	if err != nil {
		return err
	}

	conn.Close()

	return nil
}

func CreateDatabase(dbname string) {
	err := RunDatabaseCommand(fmt.Sprintf(`CREATE DATABASE %q;`, dbname))

	if err != nil {
		panic(err)
	}
}

func DropDatabase(dbname string) {
	err := RunDatabaseCommand(fmt.Sprintf(`DROP DATABASE IF EXISTS %q;`, dbname))

	if err != nil {
		panic(err)
	}
}

func GetGormTools(tx *sql.Tx) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: tx,
	}), &gorm.Config{})

	return gormDB, err
}

func Connect(params ConnectParams) {
	if instance != nil {
		return
	}

	pgConfig := postgres.Config{
		DSN:                  os.Getenv("DB_DSN"),
		PreferSimpleProtocol: true,
	}

	if params.ConfigOverride != nil {
		pgConfig = *params.ConfigOverride
	}

	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel()),
	})

	if err != nil {
		log.Info(err.Error())
		panic("failed to connect database")
	}

	instance = db

	conn, err := db.DB()
	if err != nil {
		log.WithError(err).Panic("failed to get database connection")
		panic("failed to get database connection")
	}

	conn.SetMaxIdleConns(MaxIdleConns)
	conn.SetMaxOpenConns(MaxOpenConns)
	conn.SetConnMaxLifetime(time.Hour)
}

func Close() error {
	dbConfig, err := instance.DB()
	if err != nil {
		return err
	}

	err = dbConfig.Close()
	if err != nil {
		return err
	}

	instance = nil

	return nil
}

func Get() *gorm.DB {
	return instance
}

func getLogLevel() logger.LogLevel {
	var logLevel logger.LogLevel

	level := os.Getenv(config.LogLevel)

	switch level {
	case "debug":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Silent
	}

	if os.Getenv(config.AppEnvironment) == config.Production {
		logLevel = logger.Silent
	}

	return logLevel
}

// CreateDatabaseMiddleware : to setup the database transaction middleware.
func CreateDatabaseMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			txHandle := instance.Begin()

			log.Debug("beginning database transaction")

			defer func() {
				if r := recover(); r != nil {
					txHandle.Rollback()
				}
			}()

			ctx.Set(config.TransactionGetter, txHandle)

			if slices.Contains([]int{http.StatusOK, http.StatusCreated}, ctx.Response().Status) {
				log.Debug("committing transactions")

				if err := txHandle.Commit().Error; err != nil {
					log.WithError(err).Error("trx commit error: ", err.Error())
				}
			} else {
				log.Error("rolling back transaction due to status code: ", ctx.Response().Status)
				txHandle.Rollback()
			}

			return next(ctx)
		}
	}
}
