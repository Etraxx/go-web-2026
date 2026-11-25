// GameDex : serveur web du cours (mini TP 2), qui appelle la mini API (mini TP 3)
// et affiche ses pages avec des templates (mini TP 4).
// Lancer d'abord la mini API (go run ./api depuis la racine du dépôt),
// puis depuis le dossier gamedex : go run . (ou go run . -port 9090)
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
)

// ListePage contient les données de la page liste.html.
type ListePage struct {
	Games []api.Game
	Query string
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

// messageAPIIndisponible est le message montré quand la mini API ne répond pas correctement.
const messageAPIIndisponible = "Le service de données est indisponible, réessayez plus tard"

// nombreDeJeux est le nombre de jeux dans la mini API.
const nombreDeJeux = 12

// tmpl contient tous les fichiers du dossier web/templates, lus une seule fois au démarrage.
// GameDex doit donc être lancé depuis le dossier gamedex.
var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

func main() {
	// Option -port pour choisir le port (8080 par défaut)
	port := flag.String("port", "8080", "port d'écoute")
	flag.Parse()

	// Le routeur qui associe chaque route à une fonction
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	// Toutes les autres adresses : notre page 404
	mux.HandleFunc("GET /", pageIntrouvable)

	// Les fichiers du dossier web/static (le CSS)
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Démarrage du serveur sur le port choisi
	fmt.Println("Serveur sur http://localhost:" + *port)
	err := http.ListenAndServe(":"+*port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

// accueil affiche la page d'accueil avec le template accueil.html.
func accueil(w http.ResponseWriter, r *http.Request) {
	afficherPage(w, "accueil.html", nil)
}

// listeJeux affiche tous les jeux de l'API avec le template liste.html.
func listeJeux(w http.ResponseWriter, r *http.Request) {
	// On demande les jeux à la mini API
	jeux, err := api.FetchGames()
	if err != nil {
		log.Println("FetchGames :", err)
		afficherErreur(w, http.StatusBadGateway, messageAPIIndisponible)
		return
	}

	// On lit ce que l'utilisateur a cherché (?q=...), seulement pour l'afficher
	recherche := r.URL.Query().Get("q")

	// On remplit la page avec les jeux et la recherche
	donnees := ListePage{Games: jeux, Query: recherche}
	afficherPage(w, "liste.html", donnees)
}

// detailJeu affiche la fiche d'un jeu à partir de son identifiant.
func detailJeu(w http.ResponseWriter, r *http.Request) {
	// On lit l'identifiant dans l'URL et on le transforme en nombre
	id, err := strconv.Atoi(r.PathValue("id"))

	// Pas un nombre, ou un nombre plus petit que 1 : requête invalide (400)
	if err != nil || id < 1 {
		afficherErreur(w, http.StatusBadRequest, "L'identifiant du jeu doit être un nombre plus grand que 0.")
		return
	}

	// Il n'y a que 12 jeux : au-delà, le jeu n'existe pas (404)
	if id > nombreDeJeux {
		afficherErreur(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	// On demande les jeux et les tournois à la mini API
	jeux, err := api.FetchGames()
	if err != nil {
		log.Println("FetchGames :", err)
		afficherErreur(w, http.StatusBadGateway, messageAPIIndisponible)
		return
	}
	tournois, err := api.FetchTournaments()
	if err != nil {
		log.Println("FetchTournaments :", err)
		afficherErreur(w, http.StatusBadGateway, messageAPIIndisponible)
		return
	}

	// On cherche le jeu qui a cet identifiant : s'il n'existe pas, 404
	jeu, trouve := catalog.FindGame(jeux, id)
	if !trouve {
		afficherErreur(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	// On garde seulement les tournois de ce jeu
	tournoisDuJeu := catalog.TournamentsOfGame(tournois, id)

	// On remplit la page avec le jeu et ses tournois
	donnees := JeuPage{Game: jeu, Tournaments: tournoisDuJeu}
	afficherPage(w, "jeu.html", donnees)
}

// pageIntrouvable répond 404 pour toutes les adresses inconnues.
func pageIntrouvable(w http.ResponseWriter, r *http.Request) {
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
