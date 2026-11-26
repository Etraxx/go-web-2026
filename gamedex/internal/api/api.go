// Package api parle à la mini API : il récupère les jeux et les tournois en JSON.
// Il ne sait rien du HTML ni du navigateur.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// baseURL est l'adresse de la mini API.
const baseURL = "http://localhost:8081"

// client fait les appels à l'API : il abandonne au bout de 3 secondes.
var client = &http.Client{Timeout: 3 * time.Second}

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

// SearchGames demande à la mini API les jeux filtrés par nom (q) et par genre.
// C'est la mini API qui filtre ; elle ignore les paramètres vides (sans filtre : les 12 jeux).
func SearchGames(q, genre string) ([]Game, error) {
	// url.Values encode les paramètres (les accents, les espaces...)
	params := url.Values{}
	params.Set("q", q)
	params.Set("genre", genre)

	var jeux []Game
	err := lireJSON("/api/games?"+params.Encode(), &jeux)
	return jeux, err
}

// lireJSON appelle la mini API sur le chemin donné et range le JSON reçu dans donnees.
// Elle n'affiche rien : en cas de problème, elle renvoie une erreur.
func lireJSON(chemin string, donnees any) error {
	// L'appel HTTP : erreur si l'API est éteinte ou trop lente
	resp, err := client.Get(baseURL + chemin)
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
