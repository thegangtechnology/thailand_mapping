package vaults

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"go-template/config"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/vaults_test/$GOFILE -package=vaults_test

type Authenticator interface {
	GetJWSFromRequest(req *http.Request) (string, error)
	GetKeySet() jwk.Set
	ValidateToken(token jwt.Token, googleProjectID string) error
	ParseWithoutVerify(jws string) (jwt.Token, error)
	ParseWithKeySet(jws string) (jwt.Token, error)
}

type AuthenticatorImpl struct {
}

func (j AuthenticatorImpl) GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return "", config.ErrMissingAuthHeader
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", config.ErrMalformedAuthHeader
	}

	return strings.TrimPrefix(authHdr, prefix), nil
}

func (j AuthenticatorImpl) GetKeySet() jwk.Set {
	resp, err := http.Get(os.Getenv(config.JwtURL))
	defer func() {
		resp.Body.Close()
	}()

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)

	if err != nil {
		panic(err)
	}

	keySet := jwk.NewSet()

	for k, v := range data {
		pubKey, err := jwk.New(bytesToPublicKey([]byte(v)).PublicKey)
		if err != nil {
			logrus.WithError(err).Panic("failed to create JWK")
		}

		err = pubKey.Set(jwk.AlgorithmKey, jwa.RS256)
		if err != nil {
			return nil
		}

		err = pubKey.Set(jwk.KeyIDKey, k)
		if err != nil {
			return nil
		}

		keySet.Add(pubKey)
	}

	return keySet
}

func NewAuthenticator() Authenticator {
	return AuthenticatorImpl{}
}

func bytesToPublicKey(pub []byte) *x509.Certificate {
	block, _ := pem.Decode(pub)
	if block == nil {
		panic("failed to parse PEM block containing the templates key")
	}

	publicKey, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded templates key: " + err.Error())
	}

	return publicKey
}

func (j AuthenticatorImpl) ValidateToken(token jwt.Token, googleProjectID string) error {
	err := jwt.Validate(token, jwt.WithIssuer("https://securetoken.google.com/"+googleProjectID),
		jwt.WithAudience(googleProjectID))

	return err
}

func (j AuthenticatorImpl) ParseWithoutVerify(jws string) (jwt.Token, error) {
	token, err := jwt.Parse([]byte(jws))

	return token, err
}

func (j AuthenticatorImpl) ParseWithKeySet(jws string) (jwt.Token, error) {
	token, parseErr := jwt.Parse([]byte(jws), jwt.WithKeySet(j.GetKeySet()))

	return token, parseErr
}
