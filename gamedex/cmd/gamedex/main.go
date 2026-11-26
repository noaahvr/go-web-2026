package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"gamedex/internal/server"
)

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	cat, err := chargerCatalogue()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err)
		os.Exit(1)
	}

	srv, err := server.New(cat, api.FetchGames, "web")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err, "(lancez GameDex depuis le dossier gamedex/)")
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("GameDex écoute sur http://localhost:%d", *port)
	log.Fatal(http.ListenAndServe(addr, srv.Handler()))
}

func chargerCatalogue() (*catalog.Catalog, error) {
	games, err := api.FetchGames("", "")
	if err != nil {
		return nil, err
	}
	tournaments, err := api.FetchTournaments()
	if err != nil {
		return nil, err
	}
	log.Printf("catalogue chargé : %d jeux, %d tournois", len(games), len(tournaments))
	return catalog.New(games, tournaments), nil
}
