package handlers

import (
	"go-template/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

//nolint:gosec // No hardcoded credentials
const (
	BearerAuthScopes = "BearerAuth.Scopes"
	StsKeyAuthScopes = "StsKeyAuth.Scopes"
)

func (w *ServerInterfaceWrapper) Version(ctx echo.Context) error {
	version, err := w.Handler.Version()
	if err != nil {
		return returnMessage(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusOK, version)
}

func (w *ServerInterfaceWrapper) ReadReport(ctx echo.Context) error {
	wrapAuth(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	version, err := w.Handler.ReadReport(uint(id))
	if err != nil {
		return returnMessage(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, version)
}

func (w *ServerInterfaceWrapper) CreateReport(ctx echo.Context) error {
	wrapAuth(ctx)

	var input models.CreateProductReportInput

	err := ctx.Bind(&input)
	if err != nil {
		return returnMessage(ctx, http.StatusBadRequest, err.Error())
	}

	version, err := w.Handler.StartTrx(ctx).CreateReport(input)
	if err != nil {
		return returnMessage(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, version)
}

func wrapAuth(ctx echo.Context) {
	ctx.Set(StsKeyAuthScopes, []string{""})
	ctx.Set(BearerAuthScopes, []string{""})
}
