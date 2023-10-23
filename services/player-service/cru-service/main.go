package main

import (
	"cru-service/routes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	routes.PlayerRoute(router)
	http.Handle("/", router)
	fmt.Println("Server running on 8123")
	http.ListenAndServe(":8123", nil)
}
