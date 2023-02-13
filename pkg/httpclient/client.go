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

func makeHttpBasicAuthReq(username, password string, req *http.Request, log logger.Logger) ([]byte, int) {
	req.Header.Set("Authorization", "Basic "+basicAuth(username, password))

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode
}

func makeHttpApiKeyReq(apiKeyHeaderName, apiKey string, req *http.Request, log logger.Logger) ([]byte, int) {
	if apiKeyHeaderName != "" {
		req.Header.Set(apiKeyHeaderName, apiKey)
	} else {
		// if there is no accompanying header name, assume it is the Authorization header that needs to be sent
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode
}

func makeSignedHttpDuoCall(req *http.Request, log logger.Logger) ([]byte, int) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, 500
	}
	return bodyBytes, resp.StatusCode
}

func setOauth2Header(newReq *http.Request, tokenMap map[string]string) {
	newReq.Header.Set("Authorization", "Bearer "+tokenMap["access_token"])
	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Accept", "application/json")

	newReq.Header.Set("Accept-Encoding", "identity")

	return
}

func Oauth2ApiRequest(headers map[string]string, method, url string, data interface{}, tokenMap map[string]string, log logger.Logger) ([]byte, int) {
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
			return nil, 500
		}
		body = bytes.NewBuffer(requestBody)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Error(err)
		return nil, 500
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	setOauth2Header(request, tokenMap)
	response, err := client.Do(request)

	fmt.Println(response.StatusCode)
	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return nil, 500
	}
	return respBody, response.StatusCode
}

func MakeOAuth2ApiRequest(headers map[string]string, url, method string, data interface{}, tokenMap map[string]string, log logger.Logger) (interface{}, int, error) {
	var response interface{}

	res, statusCode := Oauth2ApiRequest(headers, method, url, data, tokenMap, log)

	if string(res) == "" {
		return nil, 500, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err := json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return nil, 500, err
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
		req.Header.Set(k, v)
	}

	var response interface{}

	res, statusCode := makeHttpApiKeyReq(apiKeyHeaderName, apiKey, req, log)
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err := json.Unmarshal(res, &response)
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
		req.Header.Set(k, v)
	}

	var response interface{}

	res, statusCode := makeHttpBasicAuthReq(username, password, req, log)
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err := json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return nil, 500, err
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
			req.Header.Set(k, v)
		}
	}

	var response interface{}

	res, statusCode := makeSignedHttpDuoCall(req, log)
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err := json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return nil, 500, err
	}
	return response, statusCode, nil
}
