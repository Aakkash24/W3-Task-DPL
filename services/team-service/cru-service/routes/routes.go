package routes

import (
	"team-cru-service/controllers"

	"github.com/gorilla/mux"
)

func TeamRoute(router *mux.Router) {
	router.HandleFunc("/createTeam", controllers.CreateTeam()).Methods("POST")
	router.HandleFunc("/readTeam/{teamId}", controllers.ReadTeam()).Methods("GET")
	router.HandleFunc("/getAllTeams", controllers.GetAllTeams()).Methods("GET")
	router.HandleFunc("/updateTeamPoints", controllers.UpdateTeamPoints()).Methods("PUT")
	router.HandleFunc("/editTeam/{teamId}", controllers.EditTeam()).Methods("PUT")
}
