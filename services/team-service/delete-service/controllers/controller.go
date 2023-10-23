package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	team_delete_configs "team-delete-service/configs"
	team_delete_responses "team-delete-service/responses"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var teamCollection *mongo.Collection = team_delete_configs.GetCollection(team_delete_configs.DB, "teams")

func DeleteTeam() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		// Extract the player ID from the request URL
		params := mux.Vars(r)
		teamId := params["teamId"]

		// Convert the player ID to an ObjectID
		objID, err := primitive.ObjectIDFromHex(teamId)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := team_delete_responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Team ID", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Define the filter based on the player ID
		filter := bson.M{"id": objID}

		// Delete the player document
		result, err := teamCollection.DeleteOne(ctx, filter)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := team_delete_responses.UserResponse{Status: http.StatusInternalServerError, Message: "Server Error", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount == 0 {
			rw.WriteHeader(http.StatusNotFound)
			response := team_delete_responses.UserResponse{Status: http.StatusNotFound, Message: "Player Not Found", Data: map[string]interface{}{"message": "Player with given ID not found"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		responseMessage := fmt.Sprintf("Team with ID %s deleted", teamId)
		rw.WriteHeader(http.StatusOK)
		response := team_delete_responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"message": responseMessage}}
		json.NewEncoder(rw).Encode(response)
	}
}
