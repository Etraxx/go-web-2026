// GameDex : premier serveur web du cours (mini TP 2).
// Lancer depuis le dossier gamedex : go run . (ou go run . -port 9090)
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// nombreDeJeux est le nombre de jeux dans la mini API.
const nombreDeJeux = 12

// pageAccueil est le HTML de la page d'accueil.
const pageAccueil = `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <title>GameDex</title>
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <h1>GameDex</h1>
  <p>Bienvenue sur GameDex, le catalogue de jeux vidéo et de leurs tournois.</p>
  <p><a href="/jeux">Voir la liste des jeux</a></p>
</body>
</html>
`

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

// accueil affiche la page d'accueil en HTML, avec le CSS.
func accueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, pageAccueil)
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
