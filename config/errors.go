package config

import (
	"errors"

	"github.com/labstack/gommon/log"
)

var ErrNotFoundHTTPCode = errors.New("404")
var (
	ErrUnauthorizedHTTPCode = errors.New("401")
	ErrPortValueMissing     = errors.New("port value is missing")
	ErrDBNameValueMissing   = errors.New("dbname value is missing")
	ErrHostValueMissing     = errors.New("host value is missing")
	ErrPasswordValueMissing = errors.New("password value is missing")
	ErrUserValueMissing     = errors.New("user value is missing")
	ErrMissingAuthHeader    = errors.New("authorization header is missing")
	ErrMalformedAuthHeader  = errors.New("authorization header is malformed")
	ErrAuthInject           = errors.New("missing auth injection")
	ErrJWTAuth              = errors.New(`failed to validate token - jwt authenticator`)
	ErrUserCreate           = errors.New("cannot create user")
	ErrTokenNotFound        = errors.New("token not found in context")
	ErrNoRoles              = errors.New("no roles are given")
	ErrNoUserGiven          = errors.New("amount of users must be 2")
	ErrCannotDetectRoles    = errors.New("no roles are detected")
	ErrCannotCastGorm       = errors.New("cannot cast object to gorm")
	ErrHTTPCheckFailed      = errors.New("bad response status")
)

var (
	ErrNameNotFoundInClaims  = errors.New("name field not found in claims")
	ErrEmailNotFoundInClaims = errors.New("email field not found in claims")
	ErrNameCast              = errors.New("cannot cast name to string")
	ErrEmailCast             = errors.New("cannot cast email to string")
)

//nolint:unused // will use later
func handleError(err error) {
	switch {
	case errors.Is(err, ErrNotFoundHTTPCode):
		log.Debug("Will return 404")
	case errors.Is(err, ErrUnauthorizedHTTPCode):
		log.Debug("Will return 401")
	default:
		log.Debug("Will return 500")
	}

	log.Debug(err.Error())
}
