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
	"strings"

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
	teamSet = make(map[int]string)
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

// sanitize casing
func isNotMyTeam(myTeam, team string) bool{
	return strings.ToLower(myTeam) != strings.ToLower(team)
}

func sortDoubleHeaders(games []models.Game) ([]models.Game, error){
	/*
		cases:
			1) both games in future or past -- earlier game first
			if either game is live:
				2) if earlier game is live -- earlier game first
				3) if later game is live -- earlier game second, later game first
	*/

	// we have a double header'
	// TODO: get true first and second game
	var chronoFirst, chronoSecond models.Game
	zerothIdxGame := games[0] 
	firstIdxGame := games[1] 
	if zerothIdxGame.DoubleHeader == "Y" {
		// traditional double header, startTimeTBD = true for second game
		if firstIdxGame.Status.StartTimeTBD {
			chronoFirst = zerothIdxGame
			chronoSecond = firstIdxGame
		}else if zerothIdxGame.Status.StartTimeTBD{
			chronoFirst = firstIdxGame
			chronoSecond = zerothIdxGame
		}
	}else if zerothIdxGame.DoubleHeader == "S"{
		// split admission, compare gameDate
		zt, err := time.Parse(time.RFC3339, zerothIdxGame.GameDate)
		if err != nil{

		}
		ft, err := time.Parse(time.RFC3339, firstIdxGame.GameDate)
		if err != nil{

		}

		if ft.After(zt){
			chronoFirst = zerothIdxGame
			chronoSecond = firstIdxGame
		}else if zt.After(ft){
			chronoFirst = firstIdxGame
			chronoSecond = zerothIdxGame
		}
	}

	// ok we have our true chronological ordering
	if chronoSecond.Status.AbstractGameCode == "L"{
		return []models.Game{chronoSecond, chronoFirst}, nil
	}

	return []models.Game{chronoFirst, chronoSecond}, nil
}

func GetSchedule(c *gin.Context) {
	var myTeam string
	date := c.Query("date")
	teamId := c.Query("teamId")
	id, err:= strconv.Atoi(teamId)
	if err != nil{
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// validate requested team ID exists
	setLock.RLock()
	myTeam = teamSet[id]
	if myTeam == "" {
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


	i := 0
	myTeamsGames := make([]models.Game, 0)
	for _, x := range schedResp.Dates[0].Games {
		// in place rewrite
		if isNotMyTeam(myTeam, x.Teams.Home.Team.Name) && isNotMyTeam(myTeam, x.Teams.Away.Team.Name){
			// copy and increment index
			schedResp.Dates[0].Games[i] = x
			i++
			continue
		}

		// build slice of games w/ requested team
		myTeamsGames = append(myTeamsGames, x)
	}
	schedResp.Dates[0].Games = schedResp.Dates[0].Games[:i]

	// my team has a double header
	if len(myTeamsGames) == 2 {
		// sortDoubleHeaders(myTeamsGames)
	}

	//build response with myTeamsGames first
	tmp := schedResp.Dates[0].Games
	schedResp.Dates[0].Games = make([]models.Game, 0)
	schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, myTeamsGames...)
	schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, tmp...)

	fmt.Printf("MY TEAMS GAMES: %+v\n", myTeamsGames)
	fmt.Printf("REST OF TEAMS GAMES: %+v\n",schedResp.Dates[0].Games)
	// c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}
