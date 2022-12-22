package app

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getPageRange(r *http.Request, numPages int) (int, int, error) {
	var err error
	pageStart := 1
	pageEnd := 1

	if r.FormValue("pageStart") != "" {
		pageStart, err = strconv.Atoi(r.FormValue("pageStart"))
		pageEnd = pageStart
		if err != nil {
			return 0, 0, err
		}
	}

	if r.FormValue("pageEnd") != "" {
		pageEnd, err = strconv.Atoi(r.FormValue("pageEnd"))
		if err != nil {
			return 0, 0, err
		}
	}

	if pageStart > numPages || pageEnd > numPages || pageStart < 1 {
		return 0, 0, errors.New("invalid pageStart or pageEnd value")
	}

	if r.FormValue("allPages") != "" && r.FormValue("allPages") == "true" {
		pageStart = 1
		pageEnd = numPages
	}

	return pageStart, pageEnd, nil
}

var spaceReplacer *strings.Replacer = strings.NewReplacer("+", "%20")

func canonParams(params url.Values) string {
	// Values must be in sorted order
	for key, val := range params {
		sort.Strings(val)
		params[key] = val
	}
	// Encode will place Keys in sorted order
	ordered_params := params.Encode()
	// Encoder turns spaces into +, but we need %XX escaping
	return spaceReplacer.Replace(ordered_params)
}

func canonicalize(method string,
	host string,
	uri string,
	params url.Values,
	date string) string {
	var canon [5]string
	canon[0] = date
	canon[1] = strings.ToUpper(method)
	canon[2] = strings.ToLower(host)
	canon[3] = uri
	canon[4] = canonParams(params)
	return strings.Join(canon[:], "\n")
}

func sign(ikey string,
	skey string,
	method string,
	host string,
	uri string,
	date string,
	params url.Values) string {
	canon := canonicalize(method, host, uri, params, date)
	mac := hmac.New(sha512.New, []byte(skey))
	mac.Write([]byte(canon))
	sig := hex.EncodeToString(mac.Sum(nil))
	auth := ikey + ":" + sig
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
