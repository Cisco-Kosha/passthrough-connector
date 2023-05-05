package app

import (
	"github.com/gorilla/mux"
	"github.com/kosha/passthrough-connector/pkg/config"
	"github.com/kosha/passthrough-connector/pkg/logger"
	"os"
	"testing"
)

var (
	r *mux.Router
	logging logger.Logger
	redisAddr string
)

func init() {
	r = mux.NewRouter().StrictSlash(true)
	logging = logger.New("app", "passthrough-connector-test")
}

func TestMain(m *testing.M) {

	//Run tests
	code := m.Run()
	os.Exit(code)

}

func TestNoAuth(t *testing.T) {
	t.Setenv("SERVER_URL", "http://httpbin.org")
	t.Setenv("AUTH_TYPE", "NONE")

	cfg := config.Get()
	a := App {
		r,
		logging,
		cfg,
	}
	a.TestCommonMiddlewareNoAuth(t)
}

func TestApiKeyAuthCustomHeader(t *testing.T) {
	t.Setenv("SERVER_URL", "http://httpbin.org")
	t.Setenv("AUTH_TYPE", "API_KEY")
	t.Setenv("API_KEY", "12345678")
	t.Setenv("API_KEY_HEADER_NAME", "x-test-header")

	cfg := config.Get()
	a := App {
		r,
		logging,
		cfg,
	}
	a.TestCommonMiddlewareApiKeyCustomHeader(t)
}

func TestApiKeyAuthDefaultHeader(t *testing.T) {
	t.Setenv("SERVER_URL", "http://httpbin.org")
	t.Setenv("AUTH_TYPE", "API_KEY")
	t.Setenv("API_KEY", "12345678")

	cfg := config.Get()
	a := App {
		r,
		logging,
		cfg,
	}
	a.TestCommonMiddlewareApiKeyDefaultHeader(t)
}


func TestBearerTokenAuth(t *testing.T) {

	t.Setenv("SERVER_URL", "http://httpbin.org")
	t.Setenv("AUTH_TYPE", "BEARER_TOKEN")
	t.Setenv("BEARER_TOKEN", "test")

	cfg := config.Get()
	a := App{
		r,
		logging,
		cfg,
	}

	a.TestCommonMiddlewareBearerToken(t)
}

func TestBasicAuth(t *testing.T) {

	t.Setenv("SERVER_URL", "http://httpbin.org")
	t.Setenv("AUTH_TYPE", "BASIC_AUTH")
	t.Setenv("USERNAME", "foo")
	t.Setenv("PASSWORD", "bar")

	cfg := config.Get()
	a := App{
		r,
		logging,
		cfg,
	}

	a.TestCommonMiddlewareBasicAuth(t)
}

func TestHMAC(t *testing.T) {

	t.Setenv("SERVER_URL", "http://httpbin.org")
	t.Setenv("AUTH_TYPE", "HMAC")
	t.Setenv("IKEY", "12345678")
	t.Setenv("SKEY", "87654321")

	cfg := config.Get()
	a := App{
		r,
		logging,
		cfg,
	}

	a.TestCommonMiddlewareHMAC(t)
}

func TestOAuth(t *testing.T) {

	t.Setenv("SERVER_URL", "http://httpbin.org")
	t.Setenv("AUTH_TYPE", "OAUTH2")
	t.Setenv("ACCESS_TOKEN", "12345678")
	t.Setenv("REFRESH_TOKEN", "87654321")
	t.Setenv("EXPIRES_AT", "93048239")

	cfg := config.Get()
	a := App{
		r,
		logging,
		cfg,
	}

	a.TestCommonMiddlewareOAuth(t)
}