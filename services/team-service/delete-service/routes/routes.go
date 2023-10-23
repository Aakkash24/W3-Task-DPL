package team_delete_routes

import (
	"team-delete-service/controllers"

	"github.com/gorilla/mux"
)

func TeamRoute(router *mux.Router) {
	router.HandleFunc("/deleteTeam/{teamId}", controllers.DeleteTeam()).Methods("DELETE")
}
