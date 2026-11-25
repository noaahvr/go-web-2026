package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur : impossible de charger les templates (lancez GameDex depuis le dossier gamedex/) :", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", ficheJeu)
	mux.HandleFunc("GET /", pageIntrouvable)

	fichiers := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fichiers))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("GameDex écoute sur http://localhost:%d", *port)
	log.Fatal(http.ListenAndServe(addr, mux))
}
