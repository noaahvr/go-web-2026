package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates *template.Template

type pageAccueil struct {
	Titre string
}

type pageJeux struct {
	Titre string
	Jeux  []Game
	Q     string
}

type pageJeu struct {
	Titre    string
	Jeu      Game
	Tournois []Tournament
}

type pageErreur struct {
	Titre   string
	Code    int
	Message string
}

func rendre(w http.ResponseWriter, code int, nom string, donnees any) {
	var tampon bytes.Buffer
	if err := templates.ExecuteTemplate(&tampon, nom, donnees); err != nil {
		log.Printf("template %s : %v", nom, err)
		http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	tampon.WriteTo(w)
}

func pageErreurHTTP(w http.ResponseWriter, code int, titre, message string) {
	rendre(w, code, "erreur.html", pageErreur{Titre: titre, Code: code, Message: message})
}

func erreurAPI(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		pageErreurHTTP(w, http.StatusNotFound, "Jeu introuvable",
			"Aucun jeu ne correspond à cet identifiant.")
	case errors.Is(err, ErrTimeout):
		log.Printf("API trop lente : %v", err)
		pageErreurHTTP(w, http.StatusGatewayTimeout, "Service trop lent",
			"La mini API met trop de temps à répondre. Réessayez dans un instant.")
	default:
		log.Printf("API indisponible : %v", err)
		pageErreurHTTP(w, http.StatusBadGateway, "Service indisponible",
			"GameDex n'arrive pas à récupérer les données de la mini API. Réessayez dans un instant.")
	}
}

func accueil(w http.ResponseWriter, r *http.Request) {
	rendre(w, http.StatusOK, "accueil.html", pageAccueil{Titre: "Accueil"})
}

func pageIntrouvable(w http.ResponseWriter, r *http.Request) {
	pageErreurHTTP(w, http.StatusNotFound, "Page introuvable",
		"Cette page n'existe pas.")
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	games, err := FetchGames()
	if err != nil {
		erreurAPI(w, err)
		return
	}
	rendre(w, http.StatusOK, "jeux.html", pageJeux{
		Titre: "Les jeux",
		Jeux:  games,
		Q:     r.URL.Query().Get("q"),
	})
}

func ficheJeu(w http.ResponseWriter, r *http.Request) {
	texte := r.PathValue("id")
	id, err := strconv.Atoi(texte)
	if err != nil {
		pageErreurHTTP(w, http.StatusBadRequest, "Identifiant invalide",
			fmt.Sprintf("« %s » n'est pas un identifiant de jeu valide : un nombre est attendu.", texte))
		return
	}

	jeu, err := FetchGame(id)
	if err != nil {
		erreurAPI(w, err)
		return
	}
	tournois, err := FetchGameTournaments(id)
	if err != nil {
		erreurAPI(w, err)
		return
	}

	rendre(w, http.StatusOK, "jeu.html", pageJeu{
		Titre:    jeu.Name,
		Jeu:      jeu,
		Tournois: tournois,
	})
}
