// GameDex : serveur web du cours (mini TP 2), qui appelle la mini API (mini TP 3)
// et affiche ses pages avec des templates (mini TP 4).
// Lancer d'abord la mini API (go run ./api depuis la racine du dépôt),
// puis depuis le dossier gamedex : go run . (ou go run . -port 9090)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Game est un jeu, tel que la mini API le renvoie en JSON.
type Game struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year"`
	Studio string   `json:"studio"`
	Genres []string `json:"genres"`
}

// Tournament est un tournoi, tel que la mini API le renvoie en JSON.
type Tournament struct {
	ID     int     `json:"id"`
	GameID int     `json:"gameId"`
	Name   string  `json:"name"`
	City   string  `json:"city"`
	Date   string  `json:"date"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

// adresseAPI est l'adresse de la mini API.
const adresseAPI = "http://localhost:8081"

// client fait les appels à l'API : il abandonne au bout de 3 secondes.
var client = &http.Client{Timeout: 3 * time.Second}

// nombreDeJeux est le nombre de jeux dans la mini API.
const nombreDeJeux = 12

// tmpl contient tous les fichiers du dossier templates, lus une seule fois au démarrage.
// GameDex doit donc être lancé depuis le dossier gamedex.
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	// Option -port pour choisir le port (8080 par défaut)
	port := flag.String("port", "8080", "port d'écoute")
	flag.Parse()

	// Le routeur qui associe chaque route à une fonction
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	// Les fichiers du dossier static (le CSS)
	fs := http.FileServer(http.Dir("static"))
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

// listeJeux affiche la liste des jeux (pour l'instant un simple texte).
func listeJeux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des jeux")
}

// detailJeu affiche la fiche d'un jeu à partir de son identifiant.
func detailJeu(w http.ResponseWriter, r *http.Request) {
	// On lit l'identifiant dans l'URL et on le transforme en nombre
	id, err := strconv.Atoi(r.PathValue("id"))

	// Pas un nombre, ou un nombre plus petit que 1 : requête invalide (400)
	if err != nil || id < 1 {
		http.Error(w, "identifiant invalide", http.StatusBadRequest)
		return
	}

	// Il n'y a que 12 jeux : au-delà, le jeu n'existe pas (404)
	if id > nombreDeJeux {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Fiche du jeu %d\n", id)
}

// FetchGames récupère la liste des jeux auprès de la mini API.
func FetchGames() ([]Game, error) {
	var jeux []Game
	err := lireJSON("/api/games", &jeux)
	return jeux, err
}

// FetchTournaments récupère la liste des tournois auprès de la mini API.
func FetchTournaments() ([]Tournament, error) {
	var tournois []Tournament
	err := lireJSON("/api/tournaments", &tournois)
	return tournois, err
}

// lireJSON appelle la mini API sur le chemin donné et range le JSON reçu dans donnees.
// Elle n'affiche rien : en cas de problème, elle renvoie une erreur.
func lireJSON(chemin string, donnees any) error {
	// L'appel HTTP : erreur si l'API est éteinte ou trop lente
	resp, err := client.Get(adresseAPI + chemin)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// err == nil ne suffit pas : il faut aussi vérifier le code de réponse
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	// On transforme le JSON reçu en données Go
	err = json.NewDecoder(resp.Body).Decode(donnees)
	if err != nil {
		return fmt.Errorf("JSON illisible : %w", err)
	}
	return nil
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
