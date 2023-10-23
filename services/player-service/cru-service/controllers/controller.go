package controllers

import (
	"context"
	cru_model "cru-service/Models"
	"cru-service/configs"
	"cru-service/responses"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var playerCollection *mongo.Collection = configs.GetCollection(configs.DB, "players")

func InitPlayer(Name string, Jno string, Age string, Role [2]string) cru_model.Player {
	return cru_model.Player{
		Id:         primitive.NewObjectID(),
		Name:       Name,
		Jno:        Jno,
		Age:        Age,
		Role:       Role,
		BatAvg:     0.0,
		StrikeRate: 0.0,
		Econ:       0.0,
		Wickets:    0,
		Matches:    0,
		Runs:       0}
}

// Define a map for mapping user input to enum values (case-insensitive)
var roleMapping = map[string]string{
	"batsman":      cru_model.RoleBatsman,
	"bowler":       cru_model.RoleBowler,
	"wicketkeeper": cru_model.RoleWicketKeeper,
}

func CreatePlayer() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		var user cru_model.Player
		defer cancel()

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "Empty Details", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		if len(user.Role) != 2 {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Roles", Data: map[string]interface{}{"message": "Roles array must contain exactly two roles"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		mappedRoles := []string{
			roleMapping[strings.ToLower(user.Role[0])],
			roleMapping[strings.ToLower(user.Role[1])],
		}

		if mappedRoles[0] == cru_model.RoleInvalid || mappedRoles[1] == cru_model.RoleInvalid {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Roles", Data: map[string]interface{}{"message": "One or both roles are invalid"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newUser := cru_model.Player{
			Name: user.Name,
			Jno:  user.Jno,
			Age:  user.Age,
			Role: [2]string(mappedRoles),
		}
		newUser = InitPlayer(newUser.Name, newUser.Jno, newUser.Age, newUser.Role)
		_, err := playerCollection.InsertOne(ctx, newUser)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		responseMessage := fmt.Sprintf("Data inserted with the index %s", newUser.Id.Hex())
		rw.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"message": responseMessage}}
		json.NewEncoder(rw).Encode(response)
	}
}

func ReadPlayer() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["playerId"]
		var user cru_model.Player
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := playerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		user.BatAvg = float64(user.BatAvg)
		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)
	}
}

func EditPlayer() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		params := mux.Vars(r)
		playerID := params["playerId"]

		// Convert the player ID to an ObjectID
		objID, err := primitive.ObjectIDFromHex(playerID)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Player ID", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		filter := bson.M{"id": objID}

		// Decode the JSON request into a map
		var requestMap map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestMap); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Request", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		normalizedRequestMap := make(map[string]interface{})
		for key, value := range requestMap {
			normalizedRequestMap[strings.ToLower(key)] = value
		}

		update := bson.M{"$set": normalizedRequestMap}

		_, err = playerCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "Server Error", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		responseMessage := fmt.Sprintf("Player with ID %s updated", playerID)
		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"message": responseMessage}}
		json.NewEncoder(rw).Encode(response)
	}
}
