package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (a *App) TestCommonMiddlewareNoAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	req.RequestURI = "/headers"
	if err != nil {
		t.Fatal(err)
	}

	a.executeTest(t, req)
}

func (a *App) TestCommonMiddlewareApiKeyCustomHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	req.RequestURI = "/headers"
	if err != nil {
		t.Fatal(err)
	}

	var responseMap map[string]interface{}

	response := a.executeTest(t, req)

	err = json.Unmarshal(response.Body.Bytes(), &responseMap)
	if err != nil {
		a.Log.Errorf("can't unmarshalling response %v, error is: ", response.Body, err)
	}

	// checks to see if the custom header has been added to the request
	val, ok := responseMap["headers"].(map[string]interface{})["X-Test-Header"].([]interface{})
	if !ok || val[0].(string) != "12345678" {
		t.Errorf("request header does not contain proper default api key header")
	}
}

func (a *App) TestCommonMiddlewareApiKeyDefaultHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	req.RequestURI = "/headers"
	if err != nil {
		t.Fatal(err)
	}

	var responseMap map[string]interface{}

	response := a.executeTest(t, req)

	err = json.Unmarshal(response.Body.Bytes(), &responseMap)
	if err != nil {
		a.Log.Errorf("can't unmarshalling response %v, error is: ", response.Body, err)
	}

	// checks to see if the default header has been added to the request
	val, ok := responseMap["headers"].(map[string]interface{})["X"].([]interface{})
	if !ok || val[0].(string) != "12345678" {
		t.Errorf("request header does not contain proper default api key header")
	}
}

func (a *App) TestCommonMiddlewareBearerToken(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	// this uri checks for a request with a bearer token attached
	req.RequestURI = "/bearer"
	if err != nil {
		t.Fatal(err)
	}
	a.executeTest(t, req)
}

func (a *App) TestCommonMiddlewareBasicAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	// this uri checks for a request with basic auth with a username of foo and password of bar
	req.RequestURI = "/basic-auth/foo/bar"
	if err != nil {
		t.Fatal(err)
	}
	a.executeTest(t, req)
}

func (a *App) TestCommonMiddlewareHMAC(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	req.RequestURI = "/headers"
	if err != nil {
		t.Fatal(err)
	}

	var responseMap map[string]interface{}

	response := a.executeTest(t, req)

	err = json.Unmarshal(response.Body.Bytes(), &responseMap)
	if err != nil {
		a.Log.Errorf("can't unmarshalling response %v, error is: ", response.Body, err)
	}

	// checks to see if the date header has been added to the request
	_, ok := responseMap["headers"].(map[string]interface{})["Date"].(string)
	if !ok {
		t.Errorf("request header does not contain proper HMAC date header")
	}
}

func (a *App) TestCommonMiddlewareOAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	// oauth also adds a bearer token so this makes sure there are no issues and the token has been added
	req.RequestURI = "/bearer"
	if err != nil {
		t.Fatal(err)
	}
	a.executeTest(t, req)
}

func (a *App) executeTest(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := a.commonMiddleware()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	return rr
}
