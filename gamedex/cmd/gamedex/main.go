// GameDex : le point de départ du programme. main assemble les morceaux :
// il charge les données une seule fois, puis lance le serveur.
// Lancer d'abord la mini API (go run ./api depuis la racine du dépôt),
// puis depuis le dossier gamedex : go run ./cmd/gamedex (ou go run ./cmd/gamedex -port 9090)
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gamedex/internal/api"
	"gamedex/internal/server"
)

func main() {
	// Option -port pour choisir le port (8080 par défaut)
	port := flag.String("port", "8080", "port d'écoute")
	flag.Parse()

	// On charge les jeux et les tournois une seule fois, au démarrage
	// Si l'API ne répond pas, GameDex s'arrête tout de suite avec un message clair
	jeux, err := api.FetchGames()
	if err != nil {
		log.Fatal("Impossible de charger les jeux : ", err)
	}
	tournois, err := api.FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de charger les tournois : ", err)
	}

	// On lance le serveur avec les données en mémoire
	srv := &server.Server{Games: jeux, Tournaments: tournois}
	fmt.Println("GameDex sur http://localhost:" + *port)
	err = http.ListenAndServe(":"+*port, srv.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
