package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
