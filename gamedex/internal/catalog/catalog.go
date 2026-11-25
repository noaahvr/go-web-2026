package catalog

import (
	"slices"

	"gamedex/internal/api"
)

type Catalog struct {
	games       []api.Game
	gamesByID   map[int]api.Game
	tournaments map[int][]api.Tournament
}

func New(games []api.Game, tournaments []api.Tournament) *Catalog {
	c := &Catalog{
		games:       slices.Clone(games),
		gamesByID:   make(map[int]api.Game, len(games)),
		tournaments: make(map[int][]api.Tournament),
	}
	for _, g := range games {
		c.gamesByID[g.ID] = g
	}
	for _, t := range tournaments {
		c.tournaments[t.GameID] = append(c.tournaments[t.GameID], t)
	}
	return c
}

func (c *Catalog) Games() []api.Game {
	return slices.Clone(c.games)
}

func (c *Catalog) Game(id int) (api.Game, bool) {
	g, ok := c.gamesByID[id]
	return g, ok
}

func (c *Catalog) TournamentsOf(gameID int) []api.Tournament {
	return slices.Clone(c.tournaments[gameID])
}
