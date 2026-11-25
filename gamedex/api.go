package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
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

var (
	ErrNotFound = errors.New("introuvable sur l'API")
	ErrTimeout  = errors.New("l'API n'a pas répondu à temps")
)

var client = &http.Client{Timeout: 3 * time.Second}

func FetchGames() ([]Game, error) {
	var games []Game
	if err := getJSON("/games", &games); err != nil {
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

func FetchGame(id int) (Game, error) {
	var game Game
	if err := getJSON("/games/"+strconv.Itoa(id), &game); err != nil {
		return Game{}, fmt.Errorf("récupération du jeu %d : %w", id, err)
	}
	return game, nil
}

func FetchGameTournaments(id int) ([]Tournament, error) {
	var tournaments []Tournament
	if err := getJSON("/games/"+strconv.Itoa(id)+"/tournaments", &tournaments); err != nil {
		return nil, fmt.Errorf("récupération des tournois du jeu %d : %w", id, err)
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
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%w (%s)", ErrNotFound, detail)
		}
		return fmt.Errorf("l'API a répondu %s", detail)
	}

	if err := json.NewDecoder(resp.Body).Decode(cible); err != nil {
		return fmt.Errorf("réponse de %s illisible : %w", adresse, err)
	}
	return nil
}
