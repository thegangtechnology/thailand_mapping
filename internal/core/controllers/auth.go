package controllers

import (
	"go-template/config"

	"github.com/labstack/echo/v4"
)

type Auth interface {
	CreateAuthMiddleware() ([]echo.MiddlewareFunc, error)
}

func (s *ServerImpl) CreateAuthMiddleware() ([]echo.MiddlewareFunc, error) {
	if s.Authenticator == nil {
		return nil, config.ErrAuthInject
	}

	return []echo.MiddlewareFunc{s.jwtInject}, nil
}

func (s *ServerImpl) jwtInject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := (*s.Authenticator).ParseJWT(ctx.Request())

		if err == nil {
			ctx.Set(config.UserKey, token)
		}

		return next(ctx)
	}
}
