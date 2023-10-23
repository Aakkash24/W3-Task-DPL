package routes

import (
	"cru-service/controllers"

	"github.com/gorilla/mux"
)

func PlayerRoute(router *mux.Router) {
	router.HandleFunc("/createPlayer", controllers.CreatePlayer()).Methods("POST")
	router.HandleFunc("/readPlayer/{playerId}", controllers.ReadPlayer()).Methods("GET")
	router.HandleFunc("/editPlayer/{playerId}", controllers.EditPlayer()).Methods("PUT")
}
