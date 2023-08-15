package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
}

func getTeamsAPIResp() (*models.TeamsResponse, error){
	res, err := http.Get(teamsAPI)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var teamsResp *models.TeamsResponse
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &teamsResp)
	if err != nil {
		return nil, err
	}

	return teamsResp, nil
}

func InitTeamIdSet() {
	setLock = new(sync.RWMutex)
	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for {
			teamsResp, err := getTeamsAPIResp()
			if err != nil{
				fmt.Printf("got err when hitting teams API: %s\n", err.Error())
				continue
			}

			createSet(*teamsResp)
			<-ticker.C
		}
	}()
}

// sanitize casing
func isNotMyTeam(myTeam, team string) bool{
	return strings.ToLower(myTeam) != strings.ToLower(team)
}

func sortDoubleHeaders(games []models.Game) ([]models.Game, error){
	var chronoFirst, chronoSecond models.Game
	zerothIdxGame := games[0] 
	firstIdxGame := games[1] 

	switch zerothIdxGame.DoubleHeader {
	case "Y":
		// traditional double header, startTimeTBD = true for second game
		if firstIdxGame.Status.StartTimeTBD {
			chronoFirst = zerothIdxGame
			chronoSecond = firstIdxGame
		}else if zerothIdxGame.Status.StartTimeTBD{
			chronoFirst = firstIdxGame
			chronoSecond = zerothIdxGame
		}
	case "S":
		// split admission, compare gameDate
		zt, err := time.Parse(time.RFC3339, zerothIdxGame.GameDate)
		if err != nil{
			return nil, err
		}
		ft, err := time.Parse(time.RFC3339, firstIdxGame.GameDate)
		if err != nil{
			return nil, err
		}

		chronoFirst = zerothIdxGame
		chronoSecond = firstIdxGame
		 if zt.After(ft){
			chronoFirst = firstIdxGame
			chronoSecond = zerothIdxGame
		}
	default:
		chronoFirst = zerothIdxGame
		chronoSecond = firstIdxGame
	}

	// if second game is live, list it first
	if chronoSecond.Status.AbstractGameCode == "L"{
		return []models.Game{chronoSecond, chronoFirst}, nil
	}

	return []models.Game{chronoFirst, chronoSecond}, nil
}

func parseQueryParameters(id int, date string) (string, error){
	// validate requested team ID exists
	setLock.RLock()
	myTeam := teamSet[id]
	setLock.RUnlock()

	if myTeam == "" {
		return "", errors.New("team not found")
	}

	// validate timestamp
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", errors.New("invalid date string")
	}

	return myTeam, nil
}

func filterTeam(myTeam string,games []models.Game) (myTeamsGames []models.Game, otherTeamsGames []models.Game){
	i := 0
	for _, x := range games {
		if isNotMyTeam(myTeam, x.Teams.Home.Team.Name) && isNotMyTeam(myTeam, x.Teams.Away.Team.Name){
			games[i] = x
			i++
			continue
		}

		// build slice of games w/ requested team
		myTeamsGames = append(myTeamsGames, x)
	}
	return myTeamsGames, games[:i]

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

	myTeam, err = parseQueryParameters(id, date)
	if err != nil{
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	res, err := http.Get(fmt.Sprintf(scheuldeAPIFmtStr, date))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	defer res.Body.Close()

	var schedResp models.ScheduleResponse
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = json.Unmarshal(b, &schedResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	myTeamsGames := make([]models.Game, 0)
	fmt.Printf("LEN OF GAMES IN SCHED RESP: %+v\n", len(schedResp.Dates[0].Games))
	
	// filter myTeam games out of schedule response payload into standalone slices
	myTeamsGames, schedResp.Dates[0].Games = filterTeam(myTeam, schedResp.Dates[0].Games)

	//build ordered response payload
	tmp := schedResp.Dates[0].Games
	schedResp.Dates[0].Games = make([]models.Game, 0)	
	if len(myTeamsGames) == 2 {
		dhGames, err := sortDoubleHeaders(myTeamsGames)
		if err != nil{
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, dhGames...)
	}else{
		schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, myTeamsGames...)
	}
	schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, tmp...)
	fmt.Printf("LEN OF GAMES IN SCHED RESP served to client: %+v\n", len(schedResp.Dates[0].Games))
	c.JSON(http.StatusOK, schedResp)
}
