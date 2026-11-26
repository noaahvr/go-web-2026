package server

import (
	"cmp"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"gamedex/internal/api"
	"gamedex/internal/store"
)

const pseudoMax = 20

type pageAccueil struct {
	Titre string
}

type ligneJeu struct {
	api.Game
	Likes int
}

type pageJeux struct {
	Titre  string
	Jeux   []ligneJeu
	Q      string
	Genre  string
	Genres []string
	Filtre bool
}

type ligneTournoi struct {
	api.Tournament
	Inscrits []string
}

type pageJeu struct {
	Titre    string
	Jeu      api.Game
	Likes    int
	Tournois []ligneTournoi
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

func (s *Server) erreurAPI(w http.ResponseWriter, err error) {
	if errors.Is(err, api.ErrTimeout) {
		log.Printf("API trop lente : %v", err)
		s.pageErreurHTTP(w, http.StatusGatewayTimeout, "Service trop lent",
			"La mini API met trop de temps à répondre. Réessayez dans un instant.")
		return
	}
	log.Printf("API indisponible : %v", err)
	s.pageErreurHTTP(w, http.StatusBadGateway, "Service indisponible",
		"GameDex n'arrive pas à récupérer les données de la mini API. Réessayez dans un instant.")
}

func (s *Server) idChemin(w http.ResponseWriter, r *http.Request, objet string) (int, bool) {
	texte := r.PathValue("id")
	id, err := strconv.Atoi(texte)
	if err != nil {
		s.pageErreurHTTP(w, http.StatusBadRequest, "Identifiant invalide",
			fmt.Sprintf("« %s » n'est pas un identifiant de %s valide : un nombre est attendu.", texte, objet))
		return 0, false
	}
	return id, true
}

func (s *Server) listeJeux(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	genre := r.URL.Query().Get("genre")

	games, err := s.search(q, genre)
	if err != nil {
		s.erreurAPI(w, err)
		return
	}

	lignes := make([]ligneJeu, len(games))
	for i, g := range games {
		lignes[i] = ligneJeu{Game: g, Likes: store.Likes(g.ID)}
	}
	slices.SortStableFunc(lignes, func(a, b ligneJeu) int {
		return cmp.Compare(b.Likes, a.Likes)
	})

	s.rendre(w, http.StatusOK, "jeux.html", pageJeux{
		Titre:  "Les jeux",
		Jeux:   lignes,
		Q:      q,
		Genre:  genre,
		Genres: s.catalog.Genres(),
		Filtre: q != "" || genre != "",
	})
}

func (s *Server) ficheJeu(w http.ResponseWriter, r *http.Request) {
	id, ok := s.idChemin(w, r, "jeu")
	if !ok {
		return
	}

	jeu, ok := s.catalog.Game(id)
	if !ok {
		s.pageErreurHTTP(w, http.StatusNotFound, "Jeu introuvable",
			"Aucun jeu ne correspond à cet identifiant.")
		return
	}

	tournois := s.catalog.TournamentsOf(id)
	lignes := make([]ligneTournoi, len(tournois))
	for i, t := range tournois {
		lignes[i] = ligneTournoi{Tournament: t, Inscrits: store.Registrations(t.ID)}
	}

	s.rendre(w, http.StatusOK, "jeu.html", pageJeu{
		Titre:    jeu.Name,
		Jeu:      jeu,
		Likes:    store.Likes(id),
		Tournois: lignes,
	})
}

func (s *Server) aimerJeu(w http.ResponseWriter, r *http.Request) {
	id, ok := s.idChemin(w, r, "jeu")
	if !ok {
		return
	}
	if _, ok := s.catalog.Game(id); !ok {
		s.pageErreurHTTP(w, http.StatusNotFound, "Jeu introuvable",
			"Aucun jeu ne correspond à cet identifiant.")
		return
	}

	store.AddLike(id)

	http.Redirect(w, r, "/jeux/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (s *Server) inscrireTournoi(w http.ResponseWriter, r *http.Request) {
	id, ok := s.idChemin(w, r, "tournoi")
	if !ok {
		return
	}
	tournoi, ok := s.catalog.Tournament(id)
	if !ok {
		s.pageErreurHTTP(w, http.StatusNotFound, "Tournoi introuvable",
			"Aucun tournoi ne correspond à cet identifiant.")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	pseudo := strings.TrimSpace(r.PostFormValue("pseudo"))
	if !pseudoValide(pseudo) {
		s.pageErreurHTTP(w, http.StatusBadRequest, "Pseudo invalide",
			fmt.Sprintf("Le pseudo doit faire entre 1 et %d caractères, sans caractère de contrôle.", pseudoMax))
		return
	}

	if err := store.AddRegistration(id, pseudo); err != nil {
		s.pageErreurHTTP(w, http.StatusConflict, "Déjà inscrit",
			fmt.Sprintf("Le pseudo « %s » est déjà inscrit à ce tournoi.", pseudo))
		return
	}

	http.Redirect(w, r, "/jeux/"+strconv.Itoa(tournoi.GameID)+"#tournois", http.StatusSeeOther)
}

func pseudoValide(p string) bool {
	n := utf8.RuneCountInString(p)
	if !utf8.ValidString(p) || n < 1 || n > pseudoMax {
		return false
	}
	for _, c := range p {
		if unicode.IsControl(c) {
			return false
		}
	}
	return true
}
