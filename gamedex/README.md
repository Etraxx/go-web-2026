# GameDex · votre travail

C'est ici que vous écrivez GameDex, mini TP après mini TP.

Au mini TP 2, placez-vous dans ce dossier et créez votre module :

```
cd gamedex
go mod init gamedex
```

Depuis le mini TP 5, on lance le serveur depuis ce dossier avec `go run ./cmd/gamedex`.

À la fin de chaque mini TP : commit, puis push sur votre fork.

## Organisation du dossier (mini TP 5)

```
gamedex/
├── go.mod                   module « gamedex »
├── cmd/gamedex/main.go      charge les données une fois, lit -port, lance le serveur
├── internal/
│   ├── api/api.go           parle à la mini API : structs, client à 3 s, FetchGames, FetchTournaments, SearchGames
│   ├── catalog/catalog.go   les règles : FindGame, FindTournament, TournamentsOfGame, Genres, SortByLikes
│   ├── server/server.go     routes, handlers, templates, pages d'erreur
│   └── store/               données gardées en mémoire (map + mutex)
│       ├── likes.go         les likes des jeux
│       └── inscriptions.go  les pseudos inscrits à chaque tournoi
└── web/
    ├── templates/           header, footer, accueil, liste, jeu, erreur
    └── static/css/style.css
```

Qui a le droit de connaître qui : `cmd` utilise `api` et `server`, `server` utilise `catalog` et les
types de `api` et `store`, `catalog` utilise seulement les types de `api`, et `store` ne connaît que
`sync`. `catalog` ne connaît ni HTTP, ni HTML, ni `store`. Après le démarrage, le seul appel à la mini
API est la recherche de `/jeux` (`api.SearchGames`).

## Lancer GameDex

GameDex charge les jeux et les tournois **une seule fois, au démarrage**. Il faut donc d'abord lancer
la mini API, dans un premier terminal, depuis la racine du dépôt :

```
go run ./api
```

Puis, dans un second terminal, **depuis le dossier `gamedex/`** (les templates et le CSS sont lus dans
`web/`, un chemin relatif au dossier de lancement) :

```
cd gamedex
go run ./cmd/gamedex
```

Le serveur écoute sur http://localhost:8080. Pour choisir un autre port :

```
go run ./cmd/gamedex -port 9090
```

Si la mini API ne répond pas au démarrage (éteinte, trop lente, erreur, JSON illisible), GameDex
s'arrête tout de suite avec un message comme « Impossible de charger les jeux : ... ». Une fois
démarré, il garde les jeux et les tournois en mémoire : les fiches marchent même si l'API s'arrête.
Seule la page `/jeux` interroge encore la mini API à chaque visite (c'est elle qui filtre) : si
l'API ne répond pas, cette page affiche une erreur 502.

Les **likes et les inscriptions sont gardés en mémoire** (package `internal/store`) : ils
disparaissent quand GameDex s'arrête. C'est une « fausse base de données ».

## Les routes

| Route | Ce qu'elle renvoie |
| --- | --- |
| `GET /` | La page d'accueil |
| `GET /jeux` | La liste des 12 jeux, chaque nom est un lien vers sa fiche |
| `GET /jeux?q=...&genre=...` | Les jeux filtrés **par la mini API** (`/api/games?q=&genre=`) : `q` cherche dans le nom, `genre` doit correspondre exactement ; chaque filtre est facultatif. Les jeux les plus likés sont affichés en premier |
| `GET /jeux/{id}` | La fiche du jeu, le bouton « ♥ J'aime (N) », ses tournois (nom, ville, date) avec leurs inscrits et un formulaire d'inscription |
| `POST /jeux/{id}/like` | Ajoute un like au jeu, puis redirige (303) vers `/jeux/{id}` |
| `POST /tournois/{id}/inscription` | Inscrit le pseudo (champ `pseudo`, 1 à 30 caractères) au tournoi, puis redirige (303) vers la fiche du jeu |
| `GET /static/...` | Les fichiers du dossier `web/static/`, par exemple `/static/css/style.css` |
| Toute autre adresse en GET | La page d'erreur avec le code 404 |

Les codes d'erreur, tous affichés avec `web/templates/erreur.html` :

| Cas | Code |
| --- | --- |
| Identifiant qui n'est pas un nombre ou plus petit que 1 (`/jeux/abc`, `/jeux/0`) | 400 |
| Pseudo vide ou de plus de 30 caractères | 400 |
| Jeu ou tournoi inexistant (`/jeux/999`, `POST /jeux/999/like`, `POST /tournois/9999/inscription`) ou adresse inconnue | 404 |
| Mini API qui ne répond pas pendant une recherche sur `/jeux` | 502 |

Une autre méthode que GET sur une route existante (par exemple `POST /jeux`) répond 405.
