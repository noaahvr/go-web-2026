package server

import (
	"fmt"
	"net/http"
	"strconv"

	"gamedex/internal/api"
)

type pageAccueil struct {
	Titre string
}

type pageJeux struct {
	Titre string
	Jeux  []api.Game
	Q     string
}

type pageJeu struct {
	Titre    string
	Jeu      api.Game
	Tournois []api.Tournament
}

type pageErreur struct {
	Titre   string
	Code    int
	Message string
}

func (s *Server) accueil(w http.ResponseWriter, r *http.Request) {
	s.rendre(w, http.StatusOK, "accueil.html", pageAccueil{Titre: "Accueil"})
}

func (s *Server) pageIntrouvable(w http.ResponseWriter, r *http.Request) {
	s.pageErreurHTTP(w, http.StatusNotFound, "Page introuvable",
		"Cette page n'existe pas.")
}

func (s *Server) listeJeux(w http.ResponseWriter, r *http.Request) {
	s.rendre(w, http.StatusOK, "jeux.html", pageJeux{
		Titre: "Les jeux",
		Jeux:  s.catalog.Games(),
		Q:     r.URL.Query().Get("q"),
	})
}

func (s *Server) ficheJeu(w http.ResponseWriter, r *http.Request) {
	texte := r.PathValue("id")
	id, err := strconv.Atoi(texte)
	if err != nil {
		s.pageErreurHTTP(w, http.StatusBadRequest, "Identifiant invalide",
			fmt.Sprintf("« %s » n'est pas un identifiant de jeu valide : un nombre est attendu.", texte))
		return
	}

	jeu, ok := s.catalog.Game(id)
	if !ok {
		s.pageErreurHTTP(w, http.StatusNotFound, "Jeu introuvable",
			"Aucun jeu ne correspond à cet identifiant.")
		return
	}

	s.rendre(w, http.StatusOK, "jeu.html", pageJeu{
		Titre:    jeu.Name,
		Jeu:      jeu,
		Tournois: s.catalog.TournamentsOf(id),
	})
}
