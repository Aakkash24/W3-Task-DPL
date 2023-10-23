package main

import (
	"fmt"
	"net/http"
	team_delete_routes "team-delete-service/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	team_delete_routes.TeamRoute(router)
	http.Handle("/", router)
	fmt.Println("Server running on 8126")
	http.ListenAndServe(":8126", nil)
}
