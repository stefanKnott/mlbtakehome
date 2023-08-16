package models

type SpringLeagueTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Team struct {
	SpringLeague SpringLeagueTeam `json:"springLeague,omitempty"`
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Link         string           `json:"link"`
}

type TeamsResponse struct {
	Copyright string `json:"copyright"`
	Teams     []Team `json:"teams"`
}

type LeagueRecord struct {
	Wins   uint8  `json"wins"`
	Losses uint8  `json:"losses"`
	Pct    string `json:"pct"`
}

type ScheduleTeam struct {
	LeagueRecord LeagueRecord `json:"leagueRecord"`
	Score        uint8        `json:"score"`
	Team         Team         `json:"team"`
	IsWinner     bool         `json:"isWinner"`
	SplitSquad   bool         `json:"splitSquad"`
	SeriesNumber uint8        `json:"seriesNumber"`
}

type Teams struct {
	Away ScheduleTeam `json:"away"`
	Home ScheduleTeam `json"home"`
}

type Status struct {
	AbstractGameState string `json:"abstractGameState"`
	AbstractGameCode  string `json:"abstractGameCode"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
}

type Venue struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Content struct {
	Link string `json:"link"`
}

type Game struct {
	GamePk                 int     `json:"gamePk"`
	Link                   string  `json:"string"`
	GameType               string  `json:"gameType"`
	Season                 string  `json:"season"`
	GameDate               string  `json:"gameDate"`
	OfficialDate           string  `json:"officialDate"`
	Status                 Status  `json:"status"`
	Teams                  Teams   `json:"teams"`
	Venue                  Venue   `json:"venue"`
	Content                Content `json:"content"`
	IsTie                  bool    `json:"isTie"`
	GameNumber             uint8   `json:"gameNumber"`
	PublicFacing           bool    `json:"publicFacing"`
	DoubleHeader           string  `json:"doubleHeader"`
	GamedayType            string  `json:"gamedayType"`
	Tiebreaker             string  `json:"tiebreaker"`
	CalendarEventID        string  `json:"calendarEventID"`
	SeasonDisplay          string  `json:"seasonDisplay"`
	DayNight               string  `json:"dayNight"`
	ScheduledInnings       uint8   `json:"scheduledInnings"`
	ReverseHomeAwayStatus  bool    `json:"reverseHomeAwayStatus"`
	InningBreakLength      uint16  `json:"inningBreakLength"`
	GamesInSeries          uint8   `json:"gamesInSeries"`
	SeriesGameNumber       uint8   `json:"seriesGameNumber"`
	SeriesDescription      string  `json:"seriesDescription"`
	RecordSource           string  `json:"recordSource"`
	IfNecessary            string  `json:"ifNecessary"`
	IfNecessaryDescription string  `json:"ifNecessaryDescription"`
}

type Date struct {
	Date                 string `json:"date"`
	TotalItems           uint8  `json:"totalItems"`
	TotalEvents          uint8  `json:"totalEvents"`
	TotalGames           uint8  `json:"totalGames"`
	TotalGamesInProgress uint8  `json:"totalGamesInProgress"`
	Games                []Game `json:"games"`
}

type Event struct {
	// ?
}

type ScheduleResponse struct {
	Copyright            string  `json:"copyright"`
	TotalItems           uint8   `json:"totalItems"`
	TotalEvents          uint8   `json:"totalEvents"`
	TotalGames           uint8   `json:"totalGames"`
	TotalGamesInProgress uint8   `json:"totalGamesInProgress"`
	Dates                []Date  `json:"dates"`
	Events               []Event `json:"events"`
}
