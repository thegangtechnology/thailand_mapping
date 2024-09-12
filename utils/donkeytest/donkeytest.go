package donkeytest

import (
	"bytes"
	"embed"
	"go-template/config"
	"go-template/utils/fancylogger"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

//go:embed test.json
var test embed.FS

func BuildContext(method, target string, body io.Reader) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	return ctx
}

func BuildMultipartContext(method, target string) echo.Context {
	e := echo.New()
	req := buildMultiPartRequest(method, target)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	return ctx
}

func StartTest(t *testing.T) *gomock.Controller {
	t.Helper()
	fancylogger.SetupLogger(nil, true)

	controller := gomock.NewController(t)

	t.Setenv(config.AppEnvironment, config.Test)

	return controller
}

func buildMultiPartRequest(method, target string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("file", "test.json")
	if err != nil {
		return nil
	}

	file, err := test.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	_, err = fw.Write(file)
	if err != nil {
		panic(err)
	}

	writer.Close()

	req, err := http.NewRequest(method, target, bytes.NewReader(body.Bytes()))

	if err != nil {
		return nil
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}
