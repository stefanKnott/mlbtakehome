package handlers


import (
	"time"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
	"github.com/stefanKnott/mlbtakehome/pkg/models"
)

var _ = Describe("Processing Schedule API requests", Label("Schedule"), func() {

    BeforeEach(func() {
		InitTeamIdSet()
    })

    When("We have a single admission double Header", func() {
        Context("and both games are in the future", func() {
			game1 := models.Game{
				GameDate: "2024-04-01T20:00:00Z",
				Status: models.Status{
					AbstractGameState: "Preview",
					AbstractGameCode: "P",
					CodedGameState: "S",
					DetailedState: "Scheduled",
					StatusCode: "S",
					StartTimeTBD: false,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}

			game2 := models.Game{
				GameDate: "2024-04-01T20:00:00Z",
				Status: models.Status{
					AbstractGameState: "Preview",
					AbstractGameCode: "P",
					CodedGameState: "S",
					DetailedState: "Scheduled",
					StatusCode: "S",
					StartTimeTBD: true,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}

            It("the games should be chronologically ordered", func(ctx SpecContext) {
				sorted, err := sortDoubleHeaders([]models.Game{game1, game2})
				Expect(err).To(BeNil())
				Expect(sorted[0]).To(Equal(game1))
				Expect(sorted[1]).To(Equal(game2))
            }, SpecTimeout(time.Second * 5))
        })

		Context("and both games are in the past", func() {
			game1 := models.Game{
				GameDate: "2023-04-01T20:00:00Z",
				OfficialDate: "2024-04-01",
				Status: models.Status{
					AbstractGameState: "Final",
					AbstractGameCode: "F",
					CodedGameState: "F",
					DetailedState: "Final",
					StatusCode: "F",
					StartTimeTBD: false,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}

			game2 := models.Game{
				OfficialDate: "2023-04-01",
				Status: models.Status{
					AbstractGameState: "Final",
					AbstractGameCode: "F",
					CodedGameState: "F",
					DetailedState: "Final",
					StatusCode: "F",
					StartTimeTBD: true,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}
            It("the games should be chronologically ordered", func(ctx SpecContext) {
				sorted, err := sortDoubleHeaders([]models.Game{game1, game2})
				Expect(err).To(BeNil())
				Expect(sorted[0]).To(Equal(game1))
				Expect(sorted[1]).To(Equal(game2))
            }, SpecTimeout(time.Second * 5))
        })

        Context("the first game is live", func() {
			game1 := models.Game{
				GameDate: "2023-04-01T20:00:00Z",
				OfficialDate: "2024-04-01",
				Status: models.Status{
					AbstractGameState: "Live",
					AbstractGameCode: "L",
					CodedGameState: "L",
					DetailedState: "Live",
					StatusCode: "L",
					StartTimeTBD: false,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}

			game2 := models.Game{
				OfficialDate: "2023-04-01",
				Status: models.Status{
					AbstractGameState: "Preview",
					AbstractGameCode: "P",
					CodedGameState: "S",
					DetailedState: "Scheduled",
					StatusCode: "S",
					StartTimeTBD: true,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}
            It("the games should be chronologically ordered", func(ctx SpecContext) {
				sorted, err := sortDoubleHeaders([]models.Game{game1, game2})
				Expect(err).To(BeNil())
				Expect(sorted[0]).To(Equal(game1))
				Expect(sorted[1]).To(Equal(game2))
            }, SpecTimeout(time.Second * 5))
        })

		Context("and the second game is live", func() {
			game1 := models.Game{
				GameDate: "2024-04-01T20:00:00Z",
				OfficialDate: "2024-04-01",
				Status: models.Status{
					AbstractGameState: "Final",
					AbstractGameCode: "F",
					CodedGameState: "F",
					DetailedState: "Final",
					StatusCode: "F",
					StartTimeTBD: false,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}

			game2 := models.Game{
				OfficialDate: "2024-04-01",
				Status: models.Status{
					AbstractGameState: "Live",
					AbstractGameCode: "L",
					CodedGameState: "L",
					DetailedState: "Live",
					StatusCode: "L",
					StartTimeTBD: true,
				},
				Teams: models.Teams{
					Away: models.ScheduleTeam{
						Team: models.Team{
							ID: 0,
							Name: "testTeamAway",
						},
					},
					Home: models.ScheduleTeam{
						Team: models.Team{
							ID: 1,
							Name: "testTeamHome",
						},	
					},
				},
				DoubleHeader: "Y",
			}

            It("the second, live, game should be listed first", func(ctx SpecContext) {
				sorted, err := sortDoubleHeaders([]models.Game{game1, game2})
				Expect(err).To(BeNil())
				Expect(sorted[0]).To(Equal(game2))
				Expect(sorted[1]).To(Equal(game1))
            }, SpecTimeout(time.Second * 5))
        })
	})
})