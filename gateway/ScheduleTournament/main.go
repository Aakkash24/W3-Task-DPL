package scheduletournament

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Fixture struct {
	Team1 string
	Team2 string
}

type MatchResult struct {
	Team1  string
	Team2  string
	Winner string
}

type keyValue struct {
	TeamName string
	Points   int
}

type TournamentResult struct {
	Fixtures         []MatchResult
	PointsTable      []keyValue
	Final            MatchResult
	TournamentWinner string
}

func FixtureScheduler(teams []string) []Fixture {
	fixtures := []Fixture{}
	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			fixture := Fixture{teams[i], teams[j]}
			fixtures = append(fixtures, fixture)
		}
	}
	fmt.Println("Fixtures: ", fixtures)
	return fixtures
}

func SimulateMatch(fixture Fixture, teamPoints *map[string]int) MatchResult {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(2)
	matchwinner := ""
	if randomNumber == 0 {
		matchwinner = fixture.Team1
	} else {
		matchwinner = fixture.Team2
	}
	if matchwinner == fixture.Team1 {
		(*teamPoints)[fixture.Team2] -= 1
	} else {
		(*teamPoints)[fixture.Team1] -= 1
	}
	(*teamPoints)[matchwinner] += 3
	matchResult := MatchResult{fixture.Team1, fixture.Team2, matchwinner}
	return matchResult
}

func PrintPointsTable(teamPoints []keyValue) {
	fmt.Println("Points Table")
	for team, points := range teamPoints {
		fmt.Println(team, ":", points)
	}
}

func FindTopTwoTeams(teamPoints map[string]int) ([]string, []keyValue) {

	var keyValueSlice []keyValue

	for key, value := range teamPoints {
		keyValueSlice = append(keyValueSlice, keyValue{key, value})
	}

	sort.Slice(keyValueSlice, func(i, j int) bool {
		return keyValueSlice[i].Points > keyValueSlice[j].Points
	})

	var top2Keys []string
	for i := 0; i < 2 && i < len(keyValueSlice); i++ {
		top2Keys = append(top2Keys, keyValueSlice[i].TeamName)
	}
	PrintPointsTable(keyValueSlice)
	return top2Keys, keyValueSlice

}

func InitializeTeamPoints(teamnames []string, teamPoints *map[string]int) {
	for _, team := range teamnames {
		(*teamPoints)[team] = 0
	}
}

func SimulateTournament(teamnames []string) TournamentResult {
	fmt.Println("Simulating Tournament")
	fmt.Println("Team Names: ", teamnames)
	fixture := FixtureScheduler(teamnames)
	matchResults := []MatchResult{}
	teamPoints := make(map[string]int)
	InitializeTeamPoints(teamnames, &teamPoints)
	for _, match := range fixture {
		matchResult := SimulateMatch(match, &teamPoints)
		matchResults = append(matchResults, matchResult)
	}
	topTeams, finalPointsTable := FindTopTwoTeams(teamPoints)
	finalMatch := Fixture{topTeams[0], topTeams[1]}
	finalMatchResult := SimulateMatch(finalMatch, &teamPoints)
	tournamentResult := TournamentResult{Fixtures: matchResults, Final: finalMatchResult, PointsTable: finalPointsTable, TournamentWinner: finalMatchResult.Winner}
	return tournamentResult
}
