package internal

import (
	"go-template/config"
	"go-template/database"
	"go-template/internal/handlers"
	"go-template/internal/injection"
	"go-template/utils/fancylogger"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateServer() *echo.Echo {
	svr, errSvr := injection.Initialize()
	if errSvr != nil {
		logrus.WithError(errSvr).Fatalln("error creating server instance")
	}

	mw, err := svr.CreateAuthMiddleware()
	if err != nil {
		logrus.WithError(err).Fatalln("error creating middleware")
	}

	safeTrx := database.CreateDatabaseMiddleware()

	if err != nil {
		logrus.WithError(err).Fatalln("error creating middleware")
	}

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(safeTrx)
	e.Use(mw...)

	isProduction := os.Getenv(config.AppEnvironment) == config.Production

	fancylogger.SetupLogger(e, isProduction)

	handlers.RegisterHandlers(e, svr)

	return e
}
