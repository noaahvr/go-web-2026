package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"gamedex/internal/catalog"
)

type Server struct {
	catalog   *catalog.Catalog
	templates *template.Template
	staticDir string
}

func New(cat *catalog.Catalog, webDir string) (*Server, error) {
	tmpl, err := template.ParseGlob(filepath.Join(webDir, "templates", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("chargement des templates : %w", err)
	}
	return &Server{
		catalog:   cat,
		templates: tmpl,
		staticDir: filepath.Join(webDir, "static"),
	}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.accueil)
	mux.HandleFunc("GET /jeux", s.listeJeux)
	mux.HandleFunc("GET /jeux/{id}", s.ficheJeu)
	mux.HandleFunc("GET /", s.pageIntrouvable)

	fichiers := http.FileServer(http.Dir(s.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fichiers))

	return mux
}

func (s *Server) rendre(w http.ResponseWriter, code int, nom string, donnees any) {
	var tampon bytes.Buffer
	if err := s.templates.ExecuteTemplate(&tampon, nom, donnees); err != nil {
		log.Printf("template %s : %v", nom, err)
		http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	tampon.WriteTo(w)
}

func (s *Server) pageErreurHTTP(w http.ResponseWriter, code int, titre, message string) {
	s.rendre(w, code, "erreur.html", pageErreur{Titre: titre, Code: code, Message: message})
}
