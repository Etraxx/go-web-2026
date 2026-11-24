// GameDex : premier serveur web du cours (mini TP 2).
// Lancer depuis le dossier gamedex : go run .
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Le routeur qui associe chaque route à une fonction
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	// Démarrage du serveur
	fmt.Println("Serveur sur http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// accueil affiche un message de bienvenue.
func accueil(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur GameDex !")
}

// listeJeux affiche la liste des jeux (pour l'instant un simple texte).
func listeJeux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des jeux")
}

// detailJeu affiche la fiche d'un jeu à partir de son identifiant.
func detailJeu(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Fiche du jeu %s\n", id)
}
