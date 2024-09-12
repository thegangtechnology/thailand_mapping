package controllers

import (
	"go-template/config"
	"go-template/internal/core/services"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ServerImpl struct {
	Authenticator *services.Authenticator
	User          *services.User
	Report        *services.Report
	Product       *services.Product
	tx            *gorm.DB
}

type ServerController interface {
	Version() (string, error)
	StartTrx(ctx echo.Context) *ServerImpl
}

func (s *ServerImpl) StartTrx(ctx echo.Context) *ServerImpl {
	txHandle, ok := ctx.Get(config.TransactionGetter).(*gorm.DB)
	if !ok {
		return s
	}

	s.tx = txHandle

	return s
}

func (s *ServerImpl) Version() (string, error) {
	var version = os.Getenv("VERSION")

	if version == "" {
		version = "latest"
	}

	return version, nil
}
