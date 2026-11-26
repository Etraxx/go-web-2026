// Package store garde des données en mémoire, comme une « fausse base de données ».
// Tout disparaît quand le serveur s'arrête.
package store

import "sync"

// Le serveur traite plusieurs requêtes en même temps : le mutex garantit qu'une seule à la fois touche la map.
var (
	mu    sync.Mutex
	likes = map[int]int{}
)

// AddLike ajoute un like au jeu et renvoie le nouveau total.
func AddLike(gameID int) int {
	mu.Lock()
	defer mu.Unlock()
	likes[gameID]++
	return likes[gameID]
}

// Likes renvoie le nombre de likes d'un jeu.
func Likes(gameID int) int {
	mu.Lock()
	defer mu.Unlock()
	return likes[gameID]
}
