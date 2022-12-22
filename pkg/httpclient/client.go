package httpclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeHttpBasicAuthReq(username, password string, req *http.Request) ([]byte, int) {
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, 0
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return bodyBytes, resp.StatusCode
}

func makeHttpApiKeyReq(apiKeyHeaderName, apiKey string, req *http.Request) ([]byte, int) {
	req.Header.Add(apiKeyHeaderName, apiKey)
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, 0
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return bodyBytes, resp.StatusCode
}

func makeSignedHttpDuoCall(req *http.Request) ([]byte, int) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return bodyBytes, resp.StatusCode
}

func MakeHttpApiKeyCall(headers map[string]string, apiKeyHeaderName, apiKey, method, url string, body interface{}) (interface{}, int, error) {

	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var response interface{}

	res, statusCode := makeHttpApiKeyReq(apiKeyHeaderName, apiKey, req)
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	_ = json.Unmarshal(res, &response)
	return response, statusCode, nil
}

func MakeHttpBasicAuthCall(headers map[string]string, username, password, method, url string, body interface{}) (interface{}, int, error) {

	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var response interface{}

	res, statusCode := makeHttpBasicAuthReq(username, password, req)
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	_ = json.Unmarshal(res, &response)
	return response, statusCode, nil
}

func MakeSignedHttpDuoCall(headers map[string]string, method, host string, url string, body interface{}) (interface{}, int, error) {
	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, host+url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, host+url, nil)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var response interface{}

	res, statusCode := makeSignedHttpDuoCall(req)
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	_ = json.Unmarshal(res, &response)
	return response, statusCode, nil
}
