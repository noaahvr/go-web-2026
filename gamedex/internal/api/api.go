package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

const apiBase = "http://localhost:8081/api"

type Game struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year"`
	Studio string   `json:"studio"`
	Genres []string `json:"genres"`
}

type Tournament struct {
	ID     int     `json:"id"`
	GameID int     `json:"gameId"`
	Name   string  `json:"name"`
	City   string  `json:"city"`
	Date   string  `json:"date"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

var ErrTimeout = errors.New("l'API n'a pas répondu à temps")

var client = &http.Client{Timeout: 3 * time.Second}

func FetchGames(q, genre string) ([]Game, error) {
	params := url.Values{}
	if q != "" {
		params.Set("q", q)
	}
	if genre != "" {
		params.Set("genre", genre)
	}

	chemin := "/games"
	if len(params) > 0 {
		chemin += "?" + params.Encode()
	}

	var games []Game
	if err := getJSON(chemin, &games); err != nil {
		return nil, fmt.Errorf("récupération des jeux : %w", err)
	}
	return games, nil
}

func FetchTournaments() ([]Tournament, error) {
	var tournaments []Tournament
	if err := getJSON("/tournaments", &tournaments); err != nil {
		return nil, fmt.Errorf("récupération des tournois : %w", err)
	}
	return tournaments, nil
}

func getJSON(chemin string, cible any) error {
	adresse := apiBase + chemin

	resp, err := client.Get(adresse)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return fmt.Errorf("GET %s : %w (limite : %s)", adresse, ErrTimeout, client.Timeout)
		}
		return fmt.Errorf("impossible de joindre l'API (est-elle lancée avec « go run ./api » ?) : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail := resp.Status
		var apiErr struct {
			Error string `json:"error"`
		}
		if json.NewDecoder(resp.Body).Decode(&apiErr) == nil && apiErr.Error != "" {
			detail += " : " + apiErr.Error
		}
		return fmt.Errorf("l'API a répondu %s", detail)
	}

	if err := json.NewDecoder(resp.Body).Decode(cible); err != nil {
		return fmt.Errorf("réponse de %s illisible : %w", adresse, err)
	}
	return nil
}
