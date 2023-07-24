package httpclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kosha/passthrough-connector/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeHttpNoAuthReq(req *http.Request, log logger.Logger) ([]byte, int, error) {
	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode, err
}

func makeHttpBasicAuthReq(username, password string, req *http.Request, log logger.Logger) ([]byte, int, error) {
	req.Header.Set("Authorization", "Basic "+basicAuth(username, password))

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode, err
}

func makeHttpApiKeyReq(apiKeyHeaderName, apiKey string, req *http.Request, log logger.Logger) ([]byte, int, error) {
	if apiKeyHeaderName != "" {
		req.Header.Set(apiKeyHeaderName, apiKey)
	} else {
		// if there is no accompanying header name, assume there is no required header key value
		req.Header.Set("X", apiKey)
	}

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode, err
}

func makeHttpBearerTokenReq(bearerToken string, req *http.Request, log logger.Logger) ([]byte, int, error) {
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode, err
}

func makeSignedHttpDuoCall(req *http.Request, log logger.Logger) ([]byte, int, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	return bodyBytes, resp.StatusCode, err
}

func setOauth2Header(newReq *http.Request, tokenMap map[string]string) {
	newReq.Header.Set("Authorization", "Bearer "+tokenMap["access_token"])
	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Accept", "application/json")

	newReq.Header.Set("Accept-Encoding", "identity")

	return
}

func Oauth2ApiRequest(headers map[string]string, method, url string, data interface{}, tokenMap map[string]string, log logger.Logger) ([]byte, int, error) {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	var body io.Reader
	if data == nil {
		body = nil
	} else {
		var requestBody []byte
		requestBody, err := json.Marshal(data)
		if err != nil {
			log.Error(err)
			return nil, 500, err
		}
		body = bytes.NewBuffer(requestBody)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	setOauth2Header(request, tokenMap)
	response, err := client.Do(request)

	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	return respBody, response.StatusCode, err
}

func MakeOAuth2ApiRequest(headers map[string]string, url, method string, data interface{}, tokenMap map[string]string, log logger.Logger) (interface{}, int, error) {
	var response interface{}

	res, statusCode, err := Oauth2ApiRequest(headers, method, url, data, tokenMap, log)
	if err != nil {
		return nil, statusCode, err
	}
	if string(res) == "" {
		return nil, 500, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err = json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return nil, 500, err
	}
	return response, statusCode, nil

}

func MakeHttpNoAuthCall(headers map[string]string, method, url string, body interface{}, log logger.Logger) (interface{}, int, error) {
	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	for k, v := range headers {
		// remove user-agent header because discord doesn't like it?
		if k != "User-Agent" {
			req.Header.Add(k, v)
		}
	}
	var response interface{}

	res, statusCode, err := makeHttpNoAuthReq(req, log)
	if err != nil {
		return nil, statusCode, err
	}
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err = json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return string(res), 200, nil
	}
	return response, statusCode, nil

}

func MakeHttpApiKeyCall(headers map[string]string, apiKeyHeaderName, apiKey, method, url string, body interface{}, log logger.Logger) (interface{}, int, error) {

	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	for k, v := range headers {
		// remove user-agent header because discord doesn't like it?
		if k != "User-Agent" {
			req.Header.Add(k, v)
		}
	}

	var response interface{}

	res, statusCode, err := makeHttpApiKeyReq(apiKeyHeaderName, apiKey, req, log)
	if err != nil {
		return nil, statusCode, err
	}
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err = json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return string(res), 200, nil
	}
	return response, statusCode, nil
}

func MakeHttpBearerTokenCall(headers map[string]string, bearerToken, method, url string, body interface{}, log logger.Logger) (interface{}, int, error) {

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

	res, statusCode, err := makeHttpBearerTokenReq(bearerToken, req, log)
	if err != nil {
		return nil, statusCode, err
	}
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	//Convert response body to target struct
	err = json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return nil, 500, err
	}
	return response, statusCode, nil
}

func MakeHttpBasicAuthCall(headers map[string]string, username, password, method, url string, body interface{}, log logger.Logger) (interface{}, int, error) {

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

	res, statusCode, err := makeHttpBasicAuthReq(username, password, req, log)
	if err != nil {
		return nil, statusCode, err
	}
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err = json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return string(res), 200, nil
	}
	return response, statusCode, nil
}

func MakeSignedHttpDuoCall(headers map[string]string, method, host string, url string, body interface{}, log logger.Logger) (interface{}, int, error) {
	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, host+url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, host+url, nil)
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	var response interface{}

	res, statusCode, err := makeSignedHttpDuoCall(req, log)
	if err != nil {
		return nil, statusCode, err
	}
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err = json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return string(res), 200, nil
	}
	return response, statusCode, nil
}
