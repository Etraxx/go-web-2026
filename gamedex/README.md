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
│   ├── api/api.go           parle à la mini API : structs, client à 3 s, FetchGames, FetchTournaments
│   ├── catalog/catalog.go   les règles : FindGame, TournamentsOfGame
│   └── server/server.go     routes, handlers, templates, pages d'erreur
└── web/
    ├── templates/           header, footer, accueil, liste, jeu, erreur
    └── static/css/style.css
```

Qui a le droit de connaître qui : `cmd` utilise `api` et `server`, `server` utilise `catalog` et les
types de `api`, `catalog` utilise seulement les types de `api`. `server` n'appelle jamais la mini API,
et `catalog` ne connaît ni HTTP ni HTML.

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
démarré, il garde les données en mémoire et ne dépend plus de l'API.

## Les routes

| Route | Ce qu'elle renvoie |
| --- | --- |
| `GET /` | La page d'accueil |
| `GET /jeux` | La liste des 12 jeux, chaque nom est un lien vers sa fiche |
| `GET /jeux?q=...` | La même liste, avec « Vous avez cherché : ... » au-dessus (pas encore de filtre) |
| `GET /jeux/{id}` | La fiche du jeu et ses tournois (nom, ville, date) |
| `GET /static/...` | Les fichiers du dossier `web/static/`, par exemple `/static/css/style.css` |
| Toute autre adresse en GET | La page d'erreur avec le code 404 |

Les codes d'erreur, tous affichés avec `web/templates/erreur.html` :

| Cas | Code |
| --- | --- |
| Identifiant qui n'est pas un nombre ou plus petit que 1 (`/jeux/abc`, `/jeux/0`) | 400 |
| Jeu inexistant (`/jeux/999`) ou adresse inconnue | 404 |

Une autre méthode que GET sur une route existante (par exemple `POST /jeux`) répond 405.
