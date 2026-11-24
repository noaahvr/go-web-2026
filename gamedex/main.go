package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Jeu struct {
	ID  int
	Nom string
}

var jeux = []Jeu{
	{1, "Celeste"},
	{2, "Hades"},
	{3, "Rocket League"},
	{4, "Street Fighter 6"},
	{5, "Stardew Valley"},
	{6, "League of Legends"},
	{7, "Valorant"},
	{8, "Dead Cells"},
	{9, "Trackmania"},
	{10, "Mario Kart 8 Deluxe"},
	{11, "Hollow Knight"},
	{12, "Minecraft"},
}

func trouverJeu(id int) (Jeu, bool) {
	for _, j := range jeux {
		if j.ID == id {
			return j, true
		}
	}
	return Jeu{}, false
}

const pageAccueil = `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GameDex</title>
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <main>
    <h1>Game<span>Dex</span></h1>
    <p>Le catalogue de jeux vidéo et de leurs tournois.</p>
    <p><a class="bouton" href="/jeux">Voir les jeux</a></p>
  </main>
</body>
</html>
`

func accueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, pageAccueil)
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Liste des jeux :")
	for _, j := range jeux {
		fmt.Fprintf(w, "- %d : %s\n", j.ID, j.Nom)
	}
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	texte := r.PathValue("id")
	id, err := strconv.Atoi(texte)
	if err != nil {
		http.Error(w, fmt.Sprintf("identifiant invalide : %q n'est pas un nombre", texte), http.StatusBadRequest)
		return
	}

	jeu, ok := trouverJeu(id)
	if !ok {
		http.Error(w, fmt.Sprintf("aucun jeu avec l'id %d", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Jeu n°%d : %s\n", jeu.ID, jeu.Nom)
}

func afficherCatalogue() error {
	games, err := FetchGames()
	if err != nil {
		return err
	}
	tournaments, err := FetchTournaments()
	if err != nil {
		return err
	}

	parJeu := make(map[int]int)
	for _, t := range tournaments {
		parJeu[t.GameID]++
	}

	fmt.Printf("%d jeux dans le catalogue :\n", len(games))
	for _, g := range games {
		n := parJeu[g.ID]
		mot := "tournoi"
		if n > 1 {
			mot = "tournois"
		}
		fmt.Printf("- %s (%d) : %d %s\n", g.Name, g.Year, n, mot)
	}
	return nil
}

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	if err := afficherCatalogue(); err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	fichiers := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fichiers))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("GameDex écoute sur http://localhost:%d", *port)
	log.Fatal(http.ListenAndServe(addr, mux))
}
