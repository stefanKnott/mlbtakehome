package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stefanKnott/mlbtakehome/pkg/models"

)

var teamSet map[int]string
var setLock *sync.RWMutex

const (
	teamsAPI = "https://statsapi.mlb.com/api/v1/teams?season=2021&sportId=1"
	scheuldeAPIFmtStr = "https://statsapi.mlb.com/api/v1/schedule?date=%s&sportId=1&language=en"
)


func createSet(teamsResp models.TeamsResponse) {
	setLock.Lock()
	for _, team := range teamsResp.Teams {
		teamSet[team.ID] = team.Name
	}
	setLock.Unlock()

	fmt.Printf("LEN OF SET: %+v\n", len(teamSet))
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

			var teamsResp models.TeamsResponse
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

func GetLiveness(c *gin.Context) {
	// c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

func GetReadiness(c *gin.Context) {
	// c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

func GetSchedule(c *gin.Context) {
	date := c.Query("date")
	teamId := c.Query("teamId")
	id, err:= strconv.Atoi(teamId)
	if err != nil{
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// validate requested team ID exists
	setLock.RLock()
	if teamSet[id] == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	setLock.RUnlock()

	// validate timestamp
	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// TODO hit schedules API
	res, err := http.Get(fmt.Sprintf(scheuldeAPIFmtStr, date))
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	var schedResp models.ScheduleResponse
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
	}

	err = json.Unmarshal(b, &schedResp)
	if err != nil {
		fmt.Printf("error unmarshalling response: %s\n", err)
	}

	fmt.Printf("SCHEDULE RESPONSE: %+v\n", schedResp)


	// c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}
