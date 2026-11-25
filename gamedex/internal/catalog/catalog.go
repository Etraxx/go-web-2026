// Package catalog contient les règles de GameDex : retrouver un jeu, relier les tournois à un jeu.
// Il ne connaît ni HTTP, ni HTML, ni les templates.
package catalog

import "gamedex/internal/api"

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
