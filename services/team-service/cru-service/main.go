package main

import (
	"fmt"
	"net/http"
	"team-cru-service/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	routes.TeamRoute(router)
	http.Handle("/", router)
	fmt.Println("Server running on 8125")
	http.ListenAndServe(":8125", nil)
}
