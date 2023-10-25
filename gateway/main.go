package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	scheduletournament "gateway/ScheduleTournament"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type TeamData struct {
	TeamId     string `json:"teamId"`
	TeamPoints int    `json:"teamPoints"`
}

type TeamResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Data []string `json:"data"`
	} `json:"data"`
}

func forwardRequest(url string, method string, hasPlayerIDParam bool, hasTeamIDParam bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var forwardedURL string
		if hasPlayerIDParam {
			vars := mux.Vars(r)
			playerId := vars["playerId"]
			forwardedURL = fmt.Sprintf("%s/%s", url, playerId)
		} else if hasTeamIDParam {
			vars := mux.Vars(r)
			teamId := vars["teamId"]
			forwardedURL = fmt.Sprintf("%s/%s", url, teamId)
		} else {
			forwardedURL = url
		}
		fmt.Println("Forwarding request to: " + forwardedURL)
		fmt.Println("Method: " + method)
		req, err := http.NewRequest(method, forwardedURL, r.Body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Error in request", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		for header, values := range resp.Header {
			w.Header()[header] = values
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}

func UpdateTeamPoints(teamData TeamData) error {
	reqBody, err := json.Marshal(teamData)
	if err != nil {
		fmt.Println("Error in marshalling")
		return err
	}
	req, err := http.NewRequest("GET", "http://localhost:8080/updateTeamPoints", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error making request", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func ScheduleTournament(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://192.168.29.200:8125/getAllTeams")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()

	var teamResponse TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamResponse); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}
	if teamResponse.Status != http.StatusOK {
		fmt.Println("Error in response:", teamResponse.Message)
		return
	}
	teamResp := teamResponse.Data.Data
	teamData := make(map[string]TeamData)
	var teamnames []string
	for _, team := range teamResp {
		temp := strings.Split(team, ",")
		teamnames = append(teamnames, temp[0])
		teamData[temp[0]] = TeamData{TeamId: temp[1], TeamPoints: 0}
	}
	tournamentResult := scheduletournament.SimulateTournament(teamnames)
	teamPoints := make(map[string]int)
	for _, team := range tournamentResult.PointsTable {
		teamPoints[team.TeamName] = team.Points
	}
	for _, team := range tournamentResult.PointsTable {
		teamData[team.TeamName] = TeamData{TeamId: teamData[team.TeamName].TeamId, TeamPoints: team.Points}
		UpdateTeamPoints(teamData[team.TeamName])
	}
	tournamentResultJSON, err := json.Marshal(tournamentResult)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(tournamentResultJSON)
}

func main() {
	router := mux.NewRouter()

	// Home page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Home page")
	}).Methods("GET")

	router.HandleFunc("/createPlayer", forwardRequest("http://localhost:8123/createPlayer", http.MethodPost, false, false)).Methods("GET")
	router.HandleFunc("/readPlayer/{playerId}", forwardRequest("http://localhost:8123/readPlayer", http.MethodGet, true, false)).Methods("GET")
	router.HandleFunc("/editPlayer/{playerId}", forwardRequest("http://localhost:8123/editPlayer", http.MethodPut, true, false)).Methods("GET")
	router.HandleFunc("/deletePlayer/{playerId}", forwardRequest("http://localhost:8124/deletePlayer", http.MethodDelete, true, false)).Methods("GET")

	router.HandleFunc("/createTeam", forwardRequest("http://localhost:8125/createTeam", http.MethodPost, false, false)).Methods("GET")
	router.HandleFunc("/readTeam/{teamId}", forwardRequest("http://localhost:8125/readTeam", http.MethodGet, false, true)).Methods("GET")
	router.HandleFunc("/editTeam/{teamId}", forwardRequest("http://localhost:8125/editTeam", http.MethodPut, false, true)).Methods("GET")
	router.HandleFunc("/getAllTeams", forwardRequest("http://localhost:8125/getAllTeams", http.MethodGet, false, false)).Methods("GET")
	router.HandleFunc("/updateTeamPoints", forwardRequest("http://localhost:8125/updateTeamPoints", http.MethodPut, false, false)).Methods("GET")
	router.HandleFunc("/deleteTeam/{teamId}", forwardRequest("http://localhost:8126/deleteTeam", http.MethodDelete, false, true)).Methods("GET")
	router.HandleFunc("/scheduleTournament", ScheduleTournament).Methods("GET")

	fmt.Println("Starting server on port 8080")
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
