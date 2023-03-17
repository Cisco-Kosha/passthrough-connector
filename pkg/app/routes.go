package app

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kosha/passthrough-connector/pkg/httpclient"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	ApiKey    = "API_KEY"
	BasicAuth = "BASIC_AUTH"
	HMAC      = "HMAC"
	Oauth     = "OAUTH2"
	MERAKI    = "MERAKI"
)

func (a *App) commonMiddleware() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if (*r).Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		serverUrl := a.Cfg.GetServerURL()
		requestUri := r.RequestURI
		method := r.Method
		queryParams := r.URL.Query().Encode()

		var contentTypeHeaderFound bool

		serverUrl += requestUri
		if queryParams != "" && !strings.Contains(requestUri, "?") {
			serverUrl += "?" + queryParams
		}

		var c interface{}
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&c)
		defer r.Body.Close()

		headers := make(map[string]string)
		// Loop over header names
		if len(r.Header) > 0 {
			for name, values := range r.Header {
				// Loop over all values for the name.
				if strings.ToLower(name) == "content-type" {
					contentTypeHeaderFound = true
				}
				for _, value := range values {
					if name != "" && value != "" {
						headers[name] = value
					}
				}
			}
		}
		// use application/json as default content type
		if !contentTypeHeaderFound {
			headers["Content-Type"] = "application/json; charset=utf-8"
		}

		authType := a.Cfg.GetAuthType()
		switch authType {
		case ApiKey:
			apiKeyHeaderName := a.Cfg.GetApiKeyHeaderName()
			apiKey := a.Cfg.GetApiKey()

			res, statusCode, err := httpclient.MakeHttpApiKeyCall(headers, apiKeyHeaderName, apiKey, method, serverUrl, c, a.Log)
			if err != nil {
				a.Log.Errorf("Encountered an error while making a call: %v\n", err)
				respondWithError(w, statusCode, err.Error())
				return
			}
			if res == nil {
				respondWithJSON(w, statusCode, res)
			}
			respondWithJSON(w, statusCode, res)
			return
		case BasicAuth:
			username, password := a.Cfg.GetUsernameAndPassword()

			res, statusCode, err := httpclient.MakeHttpBasicAuthCall(headers, username, password, method, serverUrl, c, a.Log)
			if err != nil {
				a.Log.Errorf("Encountered an error while making a call: %v\n", err)
				respondWithError(w, statusCode, err.Error())
				return
			}
			if res == nil {
				respondWithJSON(w, statusCode, res)
			}
			respondWithJSON(w, statusCode, res)
			return
		case MERAKI:
			apiKeyHeaderName := "X-Cisco-Meraki-API-Key"
			apiKey := a.Cfg.GetApiKey()
			res, statusCode, err := httpclient.MakeHttpApiKeyCall(headers, apiKeyHeaderName, apiKey, method, serverUrl, c, a.Log)
			if err != nil {
				a.Log.Errorf("Encountered an error while making a call: %v\n", err)
				respondWithError(w, statusCode, err.Error())
				return
			}
			if res == nil {
				respondWithJSON(w, statusCode, res)
			}
			respondWithJSON(w, statusCode, res)
			return
		case HMAC:
			ikey, skey := a.Cfg.GetDuoIKeyAndSKey()
			currentTime := time.Now().UTC().Format(time.RFC1123Z)
			headers := make(map[string]string)
			headers["Authorization"] = sign(ikey, skey, method, a.Cfg.GetServerHost(), r.URL.Path, currentTime, r.URL.Query())
			headers["Date"] = currentTime
			res, statusCode, err := httpclient.MakeSignedHttpDuoCall(headers, method, a.Cfg.GetServerURL(), r.RequestURI, c, a.Log)
			if err != nil {
				a.Log.Errorf("Encountered an error while making a call: %v\n", err)
				respondWithError(w, statusCode, err.Error())
				return
			}
			if res == nil {
				respondWithJSON(w, statusCode, res)
			}
			respondWithJSON(w, statusCode, res)
			return
		case Oauth:
			accessToken := a.Cfg.GetAccessToken()
			refreshToken := a.Cfg.GetRefreshToken()
			expiresAt := a.Cfg.GetExpiresAt()

			tokenMap := make(map[string]string)
			tokenMap["access_token"] = accessToken
			tokenMap["refresh_token"] = refreshToken
			tokenMap["expires_at"] = expiresAt

			res, statusCode, err := httpclient.MakeOAuth2ApiRequest(headers, serverUrl, method, c, tokenMap, a.Log)
			if err != nil {
				a.Log.Errorf("Encountered an error while making a call: %v\n", err)
				respondWithError(w, statusCode, err.Error())
				return
			}
			if res == nil {
				respondWithJSON(w, statusCode, res)
			}
			respondWithJSON(w, statusCode, res)
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		respondWithJSON(w, http.StatusOK, "Hello world")
	})
}

func (a *App) InitializeRoutes() {
	a.Router.PathPrefix("/").Handler(a.commonMiddleware()).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")

	// Swagger
	a.Router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)
}
