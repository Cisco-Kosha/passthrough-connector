package app

import (
	"net/http"
)

// listConnectorSpecification godoc
// @Summary Get details of the connector specification
// @Description Get all environment variables that need to be supplied
// @Tags specification
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure      400  {object} string "bad request"
// @Failure      403  {object}  string "permission denied"
// @Failure      404  {object}  string "not found"
// @Failure      500  {object}  string "internal server error"
// @Router /api/v2/specification/list [get]
func (a *App) listConnectorSpecification(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	respondWithJSON(w, http.StatusOK, map[string]string{
		"API_KEY":     "Freshservice API Key",
		"DOMAIN_NAME": "Freshservice Domain Name",
	})
}
