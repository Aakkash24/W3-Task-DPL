package main

import (
	delete_routes "delete-service/routes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	delete_routes.DeleteRoute(router)
	http.Handle("/", router)
	fmt.Println("Server running on 8124")
	http.ListenAndServe(":8124", nil)
}
