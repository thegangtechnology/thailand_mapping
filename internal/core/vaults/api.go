package vaults

import (
	"go-template/config"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/vaults_test/$GOFILE -package=vaults_test

type APIImpl struct {
	httpClient http.Client
}

type API interface {
	SendRequest(req *http.Request) ([]byte, error)
	CreateRequest(httpMethod, url string, body io.Reader) (req *http.Request, err error)
}

func (e APIImpl) CreateRequest(httpMethod, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(httpMethod, url, body)

	return
}

func (e APIImpl) SendRequest(req *http.Request) ([]byte, error) {
	res, err := e.httpClient.Do(req)

	if err != nil {
		log.WithError(err).Error("could not send request")

		return nil, err
	}

	defer func() {
		res.Body.Close()
	}()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Error("could not create request")

		return nil, err
	}

	if res.StatusCode >= config.HTTPStatusCheckFail {
		return nil, config.ErrHTTPCheckFailed
	}

	return resBytes, nil
}

func NewAPI() API {
	client := http.Client{}

	return APIImpl{
		client,
	}
}
