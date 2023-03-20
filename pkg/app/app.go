package app

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kosha/passthrough-connector/pkg/config"
	"github.com/kosha/passthrough-connector/pkg/logger"
	"log"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
	Log    logger.Logger
	Cfg    *config.Config
}

func router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}

//func commonMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Do stuff here
//		fmt.Println(r.RequestURI)
//		// Call the next handler, which can be another middleware in the chain, or the final handler.
//		respondWithJSON(w, http.StatusOK, "ok")
//		//next.ServeHTTP(w, r)
//	})
//}

// Initialize creates the necessary scaffolding of the app
func (a *App) Initialize(log logger.Logger) {

	cfg := config.Get()

	a.Cfg = cfg
	a.Log = log

	//ctx := context.Background()
	//loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	//u, _ := url.Parse("https://petstore.swagger.io/v2/swagger.json")
	//doc, _ := loader.LoadFromURI(u)
	//// Validate OAS document
	//err := doc.Validate(ctx)
	//if err != nil {
	//	a.Log.Error(err)
	//}
	//
	//// check what the security scheme that is being used
	//securitySchema := doc.Components.SecuritySchemes
	//for name, schema := range securitySchema {
	//	name
	//	schema.Value.
	//}
	a.Router = router()
}

// Run starts the app and serves on the specified addr
func (a *App) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, a.Router)
	log.Fatal(http.ListenAndServe(addr, loggedRouter))
}
