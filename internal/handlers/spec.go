package handlers

import (
	"errors"
	"fmt"
	"go-template/config"
	"go-template/internal/core/controllers"
	"go-template/internal/core/services"
	"go-template/models"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ServerInterface interface {
	controllers.ReportsController
	controllers.ServerController
	controllers.Auth
}

type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type ServerParameter struct {
	Authenticator services.Authenticator
	User          services.User
	Product       services.Product
	Report        services.Report
}

func NewServer(s *ServerParameter) ServerInterface {
	return &controllers.ServerImpl{
		Authenticator: &s.Authenticator,
		User:          &s.User,
		Product:       &s.Product,
		Report:        &s.Report,
	}
}

func RegisterHandlers(router *echo.Echo, si ServerInterface) {
	baseURL := os.Getenv(config.BaseURL)

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/version", wrapper.Version)
	router.GET(baseURL+"/reports", wrapper.ReadReport)
	router.POST(baseURL+"/reports", wrapper.CreateReport)

	router.HTTPErrorHandler = errorHandler
}

func errorHandler(err error, ctx echo.Context) {
	ctx.Logger().Error(err)

	var (
		httpError *echo.HTTPError
		errCtx    error
	)

	if errors.As(err, &httpError) {
		errCtx = returnMessage(ctx, httpError.Code, fmt.Sprintf("%v", httpError.Message))
	} else {
		errCtx = returnMessage(ctx, http.StatusInternalServerError, err.Error())
	}

	if errCtx != nil {
		logrus.WithError(errCtx).Fatal()
	}
}

func returnMessage(ctx echo.Context, code int, message string) error {
	resp := models.Message{
		Code:    int64(code),
		Message: message,
	}

	return ctx.JSON(code, resp)
}
