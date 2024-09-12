package services

import (
	"fmt"
	"go-template/config"
	"go-template/internal/core/vaults"
	"net/http"
	"os"

	"github.com/lestrrat-go/jwx/jwt"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/services_test/$GOFILE -package=services_test

type Authenticator interface {
	ParseJWT(req *http.Request) (jwt.Token, error)
}

type AuthenticatorServiceImpl struct {
	authenticator vaults.Authenticator
}

func NewAuthenticatorService(auth vaults.Authenticator) Authenticator {
	return AuthenticatorServiceImpl{
		authenticator: auth,
	}
}

func (j AuthenticatorServiceImpl) ParseJWT(req *http.Request) (jwt.Token, error) {
	jws, err := j.authenticator.GetJWSFromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("getting jws: %w", err)
	}

	jwtVerify := os.Getenv(config.JWTVerification) != config.False
	if jwtVerify {
		token, parseErr := j.authenticator.ParseWithKeySet(jws)
		if parseErr != nil {
			return nil, fmt.Errorf("parse JWS: %w", err)
		}

		googleProjectID := os.Getenv(config.GoogleProjectID)

		err = j.authenticator.ValidateToken(token, googleProjectID)
		if err != nil {
			return nil, config.ErrJWTAuth
		}

		return token, nil
	}

	token, errParse := j.authenticator.ParseWithoutVerify(jws)
	if errParse != nil {
		return nil, fmt.Errorf("parse JWS: %w", err)
	}

	return token, nil
}
