package catalog

import (
	"slices"

	"gamedex/internal/api"
)

type Catalog struct {
	gamesByID       map[int]api.Game
	tournaments     map[int][]api.Tournament
	tournamentsByID map[int]api.Tournament
	genres          []string
}

func New(games []api.Game, tournaments []api.Tournament) *Catalog {
	c := &Catalog{
		gamesByID:       make(map[int]api.Game, len(games)),
		tournaments:     make(map[int][]api.Tournament),
		tournamentsByID: make(map[int]api.Tournament, len(tournaments)),
	}

	vus := make(map[string]bool)
	for _, g := range games {
		c.gamesByID[g.ID] = g
		for _, genre := range g.Genres {
			if !vus[genre] {
				vus[genre] = true
				c.genres = append(c.genres, genre)
			}
		}
	}
	slices.Sort(c.genres)

	for _, t := range tournaments {
		c.tournaments[t.GameID] = append(c.tournaments[t.GameID], t)
		c.tournamentsByID[t.ID] = t
	}
	return c
}

func (c *Catalog) Game(id int) (api.Game, bool) {
	g, ok := c.gamesByID[id]
	return g, ok
}

func (c *Catalog) TournamentsOf(gameID int) []api.Tournament {
	return slices.Clone(c.tournaments[gameID])
}

func (c *Catalog) Tournament(id int) (api.Tournament, bool) {
	t, ok := c.tournamentsByID[id]
	return t, ok
}

func (c *Catalog) Genres() []string {
	return slices.Clone(c.genres)
}
