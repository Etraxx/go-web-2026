// Package server contient le site web : routes, handlers, templates et pages d'erreur.
// Il n'appelle jamais la mini API : il travaille avec les données déjà chargées en mémoire.
package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
)

// Server contient les jeux et les tournois chargés une seule fois au démarrage.
type Server struct {
	Games       []api.Game
	Tournaments []api.Tournament
}

// ListePage contient les données de la page liste.html.
type ListePage struct {
	Games  []api.Game // les jeux à afficher
	Genres []string   // tous les genres, pour la liste déroulante
	Genre  string     // le genre choisi
	Query  string     // ce que l'utilisateur a cherché
}

// JeuPage contient les données de la page jeu.html.
type JeuPage struct {
	Game        api.Game
	Tournaments []api.Tournament
}

// ErreurPage contient les données de la page erreur.html.
type ErreurPage struct {
	Code    int
	Message string
}

// tmpl contient tous les fichiers du dossier web/templates, lus une seule fois au démarrage.
// GameDex doit donc être lancé depuis le dossier gamedex.
var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

// Routes construit le routeur qui associe chaque route à un handler.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.accueil)
	mux.HandleFunc("GET /jeux", s.listeJeux)
	mux.HandleFunc("GET /jeux/{id}", s.detailJeu)

	// Toutes les autres adresses : notre page 404
	mux.HandleFunc("GET /", s.pageIntrouvable)

	// Les fichiers du dossier web/static (le CSS)
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	return mux
}

// accueil affiche la page d'accueil avec le template accueil.html.
func (s *Server) accueil(w http.ResponseWriter, r *http.Request) {
	afficherPage(w, "accueil.html", nil)
}

// listeJeux affiche tous les jeux avec le template liste.html.
func (s *Server) listeJeux(w http.ResponseWriter, r *http.Request) {
	// On lit la recherche et le genre dans l'URL (pour l'instant, seulement pour les afficher)
	recherche := r.URL.Query().Get("q")
	genre := r.URL.Query().Get("genre")

	// On remplit la page (les genres viennent de tous les jeux chargés au démarrage)
	donnees := ListePage{
		Games:  s.Games,
		Genres: catalog.Genres(s.Games),
		Genre:  genre,
		Query:  recherche,
	}
	afficherPage(w, "liste.html", donnees)
}

// detailJeu affiche la fiche d'un jeu à partir de son identifiant.
func (s *Server) detailJeu(w http.ResponseWriter, r *http.Request) {
	// On lit l'identifiant dans l'URL et on le transforme en nombre
	id, err := strconv.Atoi(r.PathValue("id"))

	// Pas un nombre, ou un nombre plus petit que 1 : requête invalide (400)
	if err != nil || id < 1 {
		afficherErreur(w, http.StatusBadRequest, "L'identifiant du jeu doit être un nombre plus grand que 0.")
		return
	}

	// On cherche le jeu qui a cet identifiant : s'il n'existe pas, 404
	jeu, trouve := catalog.FindGame(s.Games, id)
	if !trouve {
		afficherErreur(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	// On remplit la page avec le jeu et ses tournois
	donnees := JeuPage{Game: jeu, Tournaments: catalog.TournamentsOfGame(s.Tournaments, id)}
	afficherPage(w, "jeu.html", donnees)
}

// pageIntrouvable répond 404 pour toutes les adresses inconnues.
func (s *Server) pageIntrouvable(w http.ResponseWriter, r *http.Request) {
	afficherErreur(w, http.StatusNotFound, "Cette page n'existe pas.")
}

// afficherPage remplit le template demandé avec les données et l'envoie au navigateur.
// Si le template plante, on écrit l'erreur dans le terminal.
func afficherPage(w http.ResponseWriter, nom string, donnees any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, nom, donnees)
	if err != nil {
		log.Println("erreur du template", nom, ":", err)
	}
}

// afficherErreur envoie la page erreur.html avec le bon code HTTP.
// Le code doit être écrit AVANT la page, sinon Go envoie 200.
func afficherErreur(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	donnees := ErreurPage{Code: code, Message: message}
	err := tmpl.ExecuteTemplate(w, "erreur.html", donnees)
	if err != nil {
		log.Println("erreur du template erreur.html :", err)
	}
}
