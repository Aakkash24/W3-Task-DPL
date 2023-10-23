package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	team_cru_model "team-cru-service/Models"
	team_configs "team-cru-service/configs"
	team_responses "team-cru-service/responses"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamData struct {
	TeamId     string `json:"teamId"`
	TeamPoints int    `json:"teamPoints"`
}

var teamCollection *mongo.Collection = team_configs.GetCollection(team_configs.DB, "teams")
var mu sync.RWMutex

func InitTeam(Name string, Captain string, Players []string, HomeGround string) team_cru_model.Team {
	return team_cru_model.Team{
		Id:         primitive.NewObjectID(),
		Name:       Name,
		Captain:    Captain,
		Players:    Players,
		HomeGround: HomeGround}
}

func CreateTeam() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		var team team_cru_model.Team
		defer cancel()

		if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := team_responses.UserResponse{Status: http.StatusBadRequest, Message: "Empty Details", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newTeam := team_cru_model.Team{
			Name:       team.Name,
			Captain:    team.Captain,
			Players:    team.Players,
			HomeGround: team.HomeGround,
		}
		newTeam = InitTeam(newTeam.Name, newTeam.Captain, newTeam.Players, newTeam.HomeGround)
		_, err := teamCollection.InsertOne(ctx, newTeam)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := team_responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		responseMessage := fmt.Sprintf("Data inserted with the index %s", newTeam.Id.Hex())
		rw.WriteHeader(http.StatusCreated)
		response := team_responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"message": responseMessage}}
		json.NewEncoder(rw).Encode(response)
	}
}

func ReadTeam() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		teamId := params["teamId"]
		var team team_cru_model.Team
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(teamId)

		err := teamCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&team)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := team_responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		rw.WriteHeader(http.StatusOK)
		response := team_responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": team}}
		json.NewEncoder(rw).Encode(response)
	}
}

func GetAllTeams() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := teamCollection.Find(ctx, bson.M{})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := team_responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		var teamnames []string
		for cursor.Next(ctx) {
			var team team_cru_model.Team
			err = cursor.Decode(&team)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := team_responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}
			teamnames = append(teamnames, team.Name+","+team.Id.Hex())
		}
		rw.WriteHeader(http.StatusOK)
		response := team_responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": teamnames}}
		json.NewEncoder(rw).Encode(response)
	}
}

func UpdateTeamPoints() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		mu.Lock()
		defer mu.Unlock()
		var teamData TeamData
		if err := json.NewDecoder(r.Body).Decode(&teamData); err != nil {
			http.Error(rw, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		objID, err := primitive.ObjectIDFromHex(teamData.TeamId)
		if err != nil {
			http.Error(rw, "Invalid Team ID: "+err.Error(), http.StatusBadRequest)
			return
		}
		filter := bson.M{"id": objID}
		update := bson.M{"$set": bson.M{"points": teamData.TeamPoints}}
		_, err = teamCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(rw, "Failed to update team: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return a success response or updated team details
		rw.WriteHeader(http.StatusOK)
		responseMessage := fmt.Sprintf("Team with ID %s updated", teamData.TeamId)
		response := team_responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"message": responseMessage}}
		json.NewEncoder(rw).Encode(response)
	}
}

func EditTeam() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		params := mux.Vars(r)
		teamId := params["teamId"]
		mu.Lock()
		defer mu.Unlock()
		// Convert the player ID to an ObjectID
		objID, err := primitive.ObjectIDFromHex(teamId)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := team_responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Team ID", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		filter := bson.M{"id": objID}

		// Decode the JSON request into a map
		var requestMap map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestMap); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := team_responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid Request", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Normalize field names to lowercase
		normalizedRequestMap := make(map[string]interface{})
		for key, value := range requestMap {
			normalizedRequestMap[strings.ToLower(key)] = value
		}

		update := bson.M{"$set": normalizedRequestMap}

		_, err = teamCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := team_responses.UserResponse{Status: http.StatusInternalServerError, Message: "Server Error", Data: map[string]interface{}{"error": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		responseMessage := fmt.Sprintf("Team with ID %s updated", teamId)
		rw.WriteHeader(http.StatusOK)
		response := team_responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"message": responseMessage}}
		json.NewEncoder(rw).Encode(response)
	}
}
