package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/stefanKnott/mlbtakehome/pkg/handlers"

	"github.com/gin-gonic/gin"
)

var teamSet map[int]string
var setLock *sync.RWMutex

const (
	teamsAPI = "https://statsapi.mlb.com/api/v1/teams?season=2021&sportId=1"
)

type SpringLeagueTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Team struct {
	SpringLeague SpringLeagueTeam `json:"springLeague"`
	ID           int              `json:"id"`
	Name         string           `json:"name"`
}

type TeamsResponse struct {
	Copyright string `json:"copyright"`
	Teams     []Team `json:"teams"`
}

func createSet(teamsResp TeamsResponse) {
	setLock.Lock()
	for _, team := range teamsResp.Teams {
		teamSet[team.ID] = team.Name
	}
	setLock.Unlock()
}

func InitTeamIdSet() {
	setLock = new(sync.RWMutex)
	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for {
			teamSet = make(map[int]string)
			//getreq

			//TODO: pull this out into its own function
			res, err := http.Get(teamsAPI)
			if err != nil {
				fmt.Printf("error making http request: %s\n", err)
				os.Exit(1)
			}
			defer res.Body.Close()

			var teamsResp TeamsResponse
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("error reading response body: %s\n", err)
			}

			err = json.Unmarshal(b, &teamsResp)
			if err != nil {
				fmt.Printf("error unmarshalling response: %s\n", err)
			}

			createSet(teamsResp)
			// iter thru response
			<-ticker.C
		}
	}()
}

func main() {
	// start server
	InitTeamIdSet()
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/liveness", handlers.GetLiveness)
		v1.GET("/readiness", handlers.GetReadiness)
		v1.GET("/schedule", handlers.GetSchedule)
	}
	router.Run()
}
