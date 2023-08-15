package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stefanKnott/mlbtakehome/pkg/models"
)

var teamSet map[int]string
var setLock *sync.RWMutex

const (
	teamsAPI          = "https://statsapi.mlb.com/api/v1/teams?season=2021&sportId=1"
	scheuldeAPIFmtStr = "https://statsapi.mlb.com/api/v1/schedule?date=%s&sportId=1&language=en"
)

// structs for /schedule API responses
type ScheduleResponse struct{
	models.ScheduleResponse 
}

type ScheduleErrorResponse struct{
	Message string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func createTeamsSet(teamsResp models.TeamsResponse) {
	setLock.Lock()
	teamSet = make(map[int]string)
	for _, team := range teamsResp.Teams {
		teamSet[team.ID] = team.Name
	}
	setLock.Unlock()
}

func getTeamsAPIResp() (*models.TeamsResponse, error) {
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

	// use background goroutine to recreate team set every 30 minutes incase IDs change in MLB backend
	go func() {
		for {
			teamsResp, err := getTeamsAPIResp()
			if err != nil {
				fmt.Printf("got err when hitting teams API: %s\n", err.Error())
				continue
			}

			createTeamsSet(*teamsResp)
			<-ticker.C
		}
	}()
}

// sanitize casing
func isNotMyTeam(myTeam, team string) bool {
	return strings.ToLower(myTeam) != strings.ToLower(team)
}

func sortDoubleHeaders(games []models.Game) ([]models.Game, error) {
	var chronoFirst, chronoSecond models.Game
	zerothIdxGame := games[0]
	firstIdxGame := games[1]

	switch zerothIdxGame.DoubleHeader {
	case "Y":
		// traditional double header, startTimeTBD = true for second game
		if firstIdxGame.Status.StartTimeTBD {
			chronoFirst = zerothIdxGame
			chronoSecond = firstIdxGame
		} else if zerothIdxGame.Status.StartTimeTBD {
			chronoFirst = firstIdxGame
			chronoSecond = zerothIdxGame
		}
	case "S":
		// split admission, compare gameDate
		zt, err := time.Parse(time.RFC3339, zerothIdxGame.GameDate)
		if err != nil {
			return nil, err
		}
		ft, err := time.Parse(time.RFC3339, firstIdxGame.GameDate)
		if err != nil {
			return nil, err
		}

		chronoFirst = zerothIdxGame
		chronoSecond = firstIdxGame
		if zt.After(ft) {
			chronoFirst = firstIdxGame
			chronoSecond = zerothIdxGame
		}
	default:
		chronoFirst = zerothIdxGame
		chronoSecond = firstIdxGame
	}

	// if second game is live, list it first
	if chronoSecond.Status.AbstractGameCode == "L" {
		return []models.Game{chronoSecond, chronoFirst}, nil
	}

	return []models.Game{chronoFirst, chronoSecond}, nil
}

func parseQueryParameters(id int, date string) (string, error) {
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

// give a slice of games, filter out all games that myTeam is either home or away
// return two standalone slices representing myTeamsGames and all other games
func filterTeam(myTeam string, games []models.Game) (myTeamsGames []models.Game, otherTeamsGames []models.Game) {
	i := 0
	for _, x := range games {
		// ensure neither home nor away team is myTeam
		if isNotMyTeam(myTeam, x.Teams.Home.Team.Name) && isNotMyTeam(myTeam, x.Teams.Away.Team.Name) {
			games[i] = x
			i++
			continue
		}

		// build slice of games myTeam is either home or away
		myTeamsGames = append(myTeamsGames, x)
	}
	return myTeamsGames, games[:i]

}

// GetSchedule serves the /schedule?teamId=<id>&date=<YYYY-MM-DD> API
// which allows a client to receive a list ofgames scheduled for a specific date
// with the requested team's games ordered first
func GetSchedule(c *gin.Context) {
	var myTeam string
	date := c.Query("date")
	teamId := c.Query("teamId")
	id, err := strconv.Atoi(teamId)
	if err != nil {
		c.JSON(http.StatusBadRequest, ScheduleErrorResponse{Message: err.Error(), Timestamp: time.Now().UTC().String()})
		return
	}

	myTeam, err = parseQueryParameters(id, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ScheduleErrorResponse{Message: err.Error(), Timestamp: time.Now().UTC().String()})
		return
	}

	res, err := http.Get(fmt.Sprintf(scheuldeAPIFmtStr, date))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ScheduleErrorResponse{Message: err.Error(), Timestamp: time.Now().UTC().String()})
		return
	}
	defer res.Body.Close()

	var schedResp models.ScheduleResponse
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ScheduleErrorResponse{Message: err.Error(), Timestamp: time.Now().UTC().String()})
		return
	}

	err = json.Unmarshal(b, &schedResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ScheduleErrorResponse{Message: err.Error(), Timestamp: time.Now().UTC().String()})
		return
	}

	// unexpected response, should only contain one date
	if len(schedResp.Dates) != 1 {
		c.JSON(http.StatusInternalServerError, ScheduleErrorResponse{Message: "received invalid dates slice from schedule API", Timestamp: time.Now().UTC().String()})
		return
	}

	myTeamsGames := make([]models.Game, 0)
	// filter myTeam games out of schedule response payload into standalone slices
	myTeamsGames, schedResp.Dates[0].Games = filterTeam(myTeam, schedResp.Dates[0].Games)

	// build ordered response payload
	tmp := schedResp.Dates[0].Games
	schedResp.Dates[0].Games = make([]models.Game, 0)
	if len(myTeamsGames) == 2 {
		dhGames, err := sortDoubleHeaders(myTeamsGames)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ScheduleErrorResponse{Message: err.Error(), Timestamp: time.Now().UTC().String()})
			return
		}
		schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, dhGames...)
	} else {
		schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, myTeamsGames...)
	}
	schedResp.Dates[0].Games = append(schedResp.Dates[0].Games, tmp...)

	c.JSON(http.StatusOK, ScheduleResponse{schedResp})
}
