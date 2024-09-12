//go:build wireinject

package injection

import (
	"go-template/database"
	"go-template/internal/core/controllers"
	"go-template/internal/core/services"
	"go-template/internal/core/vaults"
	"go-template/internal/handlers"

	"github.com/google/wire"
)

func initAuthenticatorService() (services.Authenticator, error) {
	wire.Build(vaults.NewAuthenticator, services.NewAuthenticatorService)

	return services.AuthenticatorServiceImpl{}, nil
}

func initUserService() (services.User, error) {
	wire.Build(database.Get, vaults.NewUser, vaults.NewAPI, services.NewUser)

	return services.UserImpl{}, nil
}

func initReportService() (services.Report, error) {
	wire.Build(database.Get, vaults.NewReport, services.NewReport)

	return services.ReportImpl{}, nil
}

func initProductService() (services.Product, error) {
	wire.Build(database.Get, vaults.NewProduct, services.NewProduct)

	return services.ProductImpl{}, nil
}

func provideServerParameter() (*handlers.ServerParameter, error) {
	wire.Build(wire.Struct(new(handlers.ServerParameter), "*"),
		initAuthenticatorService, initUserService, initReportService, initProductService)

	return &handlers.ServerParameter{}, nil
}

func Initialize() (handlers.ServerInterface, error) {
	wire.Build(provideServerParameter, handlers.NewServer)

	return &controllers.ServerImpl{}, nil
}
