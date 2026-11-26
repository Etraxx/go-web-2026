// Package catalog contient les règles de GameDex : retrouver un jeu, relier les tournois à un jeu,
// lister les genres.
// Il ne connaît ni HTTP ni HTML : seulement les données de internal/api.
package catalog

import (
	"sort"

	"gamedex/internal/api"
)

// FindGame retrouve le jeu qui a cet identifiant.
// Le booléen vaut false si aucun jeu n'a cet identifiant.
func FindGame(games []api.Game, id int) (api.Game, bool) {
	for i := 0; i < len(games); i++ {
		if games[i].ID == id {
			return games[i], true
		}
	}
	return api.Game{}, false
}

// TournamentsOfGame renvoie les tournois du jeu gameID.
// S'il n'y en a aucun, la slice renvoyée est vide (et pas nil).
func TournamentsOfGame(tournaments []api.Tournament, gameID int) []api.Tournament {
	tournois := []api.Tournament{}
	for i := 0; i < len(tournaments); i++ {
		if tournaments[i].GameID == gameID {
			tournois = append(tournois, tournaments[i])
		}
	}
	return tournois
}

// Genres renvoie tous les genres des jeux, sans doublon, triés par ordre alphabétique.
func Genres(games []api.Game) []string {
	// La map sert à savoir si on a déjà vu un genre
	dejaVu := map[string]bool{}
	genres := []string{}
	for i := 0; i < len(games); i++ {
		for j := 0; j < len(games[i].Genres); j++ {
			genre := games[i].Genres[j]
			if !dejaVu[genre] {
				dejaVu[genre] = true
				genres = append(genres, genre)
			}
		}
	}
	sort.Strings(genres)
	return genres
}
