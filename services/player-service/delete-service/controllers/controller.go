package delete_controller

import (
	"context"
	delete_configs "delete-service/configs"
	delete_responses "delete-service/responses"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var playerCollection *mongo.Collection = delete_configs.GetCollection(delete_configs.DB, "players")

func DeletePlayer() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		// Extract the player ID from the request URL
		params := mux.Vars(r)
		playerID := params["playerId"]

		// Convert the player ID to an ObjectID
		objID, err := primitive.ObjectIDFromHex(playerID)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := delete_responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Player ID", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Define the filter based on the player ID
		filter := bson.M{"id": objID}

		// Delete the player document
		result, err := playerCollection.DeleteOne(ctx, filter)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := delete_responses.UserResponse{Status: http.StatusInternalServerError, Message: "Server Error", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount == 0 {
			rw.WriteHeader(http.StatusNotFound)
			response := delete_responses.UserResponse{Status: http.StatusNotFound, Message: "Player Not Found", Data: map[string]interface{}{"message": "Player with given ID not found"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		responseMessage := fmt.Sprintf("Player with ID %s deleted", playerID)
		rw.WriteHeader(http.StatusOK)
		response := delete_responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"message": responseMessage}}
		json.NewEncoder(rw).Encode(response)
	}
}
