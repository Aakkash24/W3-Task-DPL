package delete_routes

import (
	delete_controller "delete-service/controllers"

	"github.com/gorilla/mux"
)

func DeleteRoute(router *mux.Router) {
	router.HandleFunc("/deletePlayer/{playerId}", delete_controller.DeletePlayer()).Methods("DELETE")
}
