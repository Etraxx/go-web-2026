# Mini TP 5 : ranger GameDex

Objectif : découper mon `gamedex/main.go` du mini TP 4 en packages (`cmd/`, `internal/`, `web/`),
sans changer ce que voit l'utilisateur, et charger les données une seule fois au démarrage.

## L'arborescence finale

```
gamedex/
├── go.mod
├── README.md
├── cmd/
│   └── gamedex/
│       └── main.go          charge les données une fois, lit -port, lance le serveur
├── internal/
│   ├── api/
│   │   └── api.go           structs Game et Tournament, client à 3 s, FetchGames, FetchTournaments
│   ├── catalog/
│   │   └── catalog.go       FindGame, TournamentsOfGame
│   └── server/
│       └── server.go        Server, Routes, handlers, structs de page, templates, afficherErreur
└── web/
    ├── static/
    │   └── css/
    │       └── style.css
    └── templates/
        ├── accueil.html
        ├── erreur.html
        ├── footer.html
        ├── header.html
        ├── jeu.html
        └── liste.html
```

J'ai déplacé les fichiers avec `git mv` (templates, static et main.go) pour garder leur historique.
L'ancien `gamedex/main.go` et les dossiers `gamedex/templates` et `gamedex/static` n'existent plus.
On lance maintenant GameDex depuis `gamedex/` avec `go run ./cmd/gamedex` (ou `-port 9090`).

## Comparaison avant / après

**Avant** de toucher au code, j'ai lancé la mini API et GameDex (version mini TP 4) et j'ai
enregistré avec curl le code et le corps de 12 réponses dans un dossier temporaire, hors du dépôt :

```
01 GET / -> 200 text/html; charset=utf-8
02 GET /jeux -> 200 text/html; charset=utf-8
03 GET /jeux?q=cel -> 200 text/html; charset=utf-8
04 GET /jeux?q=%3Cscript%3Ealert(1)%3C%2Fscript%3E -> 200 text/html; charset=utf-8
05 GET /jeux/1 -> 200 text/html; charset=utf-8
06 GET /jeux/3 -> 200 text/html; charset=utf-8
07 GET /jeux/abc -> 400 text/html; charset=utf-8
08 GET /jeux/0 -> 400 text/html; charset=utf-8
09 GET /jeux/999 -> 404 text/html; charset=utf-8
10 GET /nimporte/quoi -> 404 text/html; charset=utf-8
11 GET /static/css/style.css -> 200 text/css; charset=utf-8
12 POST /jeux -> 405 text/plain; charset=utf-8
```

**Après** le rangement, j'ai relancé avec `go run ./cmd/gamedex`, refait les mêmes requêtes, puis
comparé les deux dossiers avec `diff -r avant apres` :

```
$ diff avant/codes.txt apres/codes.txt && echo "codes identiques"
codes identiques
$ diff -r avant apres && echo "corps et codes IDENTIQUES (12 réponses)"
corps et codes IDENTIQUES (12 réponses)
```

`diff` n'affiche aucune différence : mêmes codes et même HTML, octet par octet, y compris le
`<script>` échappé (`Vous avez cherché : &lt;script&gt;alert(1)&lt;/script&gt;`). J'ai aussi refait
la comparaison après chaque étape du rangement, et elle était identique à chaque fois.

## Vérification de la structure

```
$ grep -rn -e 8081 -e FetchGames -e FetchTournaments -e http.Client internal/server internal/catalog
(aucun résultat)
$ grep -rn template internal/catalog
(aucun résultat)
$ go list -f "{{.ImportPath}} -> {{.Imports}}" ./...
gamedex/cmd/gamedex -> [flag fmt gamedex/internal/api gamedex/internal/server log net/http]
gamedex/internal/api -> [encoding/json fmt net/http time]
gamedex/internal/catalog -> [gamedex/internal/api]
gamedex/internal/server -> [gamedex/internal/api gamedex/internal/catalog html/template log net/http strconv]
$ grep -n "<" cmd/gamedex/main.go
(aucun résultat)
```

- Plus aucun appel à l'API dans `server` ni dans `catalog`.
- `catalog` importe seulement `gamedex/internal/api` : ni `net/http` ni `html/template`.
- `api` n'importe pas `html/template`, et `main.go` ne contient aucun HTML.
- Les flèches vont toujours dans le même sens (cmd → server → catalog → api) : pas d'import
  circulaire. `go vet ./...` et `go build ./...` passent, et `gofmt -l .` n'affiche rien.

## Les données sont chargées une seule fois

J'ai lancé la mini API en envoyant sa sortie dans un fichier (`api.log`), puis GameDex. J'ai fait les
12 requêtes de la comparaison, plus 20 requêtes (`/jeux/1` à `/jeux/12` et 8 fois `/jeux`), donc
32 requêtes en tout. Le log de l'API :

```
2026/11/25 12:20:53 Mini API sur http://localhost:8081/api/
2026/11/25 12:20:53 Documentation sur http://localhost:8081/api/docs
2026/11/25 12:20:54 GET /api/games
2026/11/25 12:20:54 GET /api/tournaments
```

Seulement 2 requêtes, toutes les deux au démarrage. Pour comparer, la version du mini TP 4 avait fait
7 appels à l'API pour les 12 requêtes de la capture « avant » (un `/api/games` par liste, et
`/api/games` + `/api/tournaments` pour chaque fiche).

## Quand l'API est éteinte

**Au démarrage** : j'ai arrêté la mini API et lancé GameDex.

```
$ go run ./cmd/gamedex
2026/11/25 12:21:08 Impossible de charger les jeux : Get "http://localhost:8081/api/games": dial tcp [::1]:8081: connect: connection refused
exit status 1
```

GameDex s'arrête tout de suite, avec un message clair et le code de sortie 1, sans panic ni stacktrace.

**Après le démarrage** : avec l'API allumée, j'ai démarré GameDex, puis j'ai arrêté l'API. `/jeux` et
`/jeux/3` répondent toujours 200 (Rocket League et ses tournois s'affichent), et `/jeux/999` 404.

Ça a changé par rapport au mini TP 4 : avant, chaque page appelait l'API, donc une API éteinte donnait
une page 502 à chaque requête. Maintenant les jeux et les tournois sont chargés une seule fois et gardés
en mémoire dans `server.Server`. Sans données, GameDex ne peut rien afficher, donc il vaut mieux refuser
de démarrer. Une fois démarré, il ne dépend plus de l'API. J'ai donc supprimé la gestion du 502 dans les
handlers, mais j'ai gardé `afficherErreur` pour le 400 et le 404.

## Qui connaît quoi

| Package | Ce qu'il connaît | Ce qu'il ne connaît pas |
| --- | --- | --- |
| `internal/api` | L'adresse de la mini API, le client HTTP à 3 s, le JSON, les structs `Game` et `Tournament` | Le HTML, les templates, le navigateur, les routes de GameDex |
| `internal/catalog` | Les types de `api` et les règles : retrouver un jeu, garder les tournois d'un jeu | HTTP, le HTML, les templates, l'adresse de l'API |
| `internal/server` | Les routes, les handlers, les templates, les codes d'erreur, `catalog` et les types de `api` | L'adresse de l'API et comment on y va (il ne fait aucun appel) |
| `web/` | Le HTML et le CSS, et les champs des structs de page | Le Go, l'API, les règles métier |
| `cmd/gamedex` | Comment tout assembler : l'option `-port`, le chargement au démarrage, le lancement | Le HTML et les règles |

## Ce que j'en retiens

Avec ce découpage, chaque problème a sa place. Un bug d'affichage se cherche dans `web/` ou
`internal/server`, une mauvaise réponse de l'API dans `internal/api`, et une erreur de règle (le mauvais
jeu, les mauvais tournois) dans `internal/catalog`. À deux, on peut travailler en même temps, par exemple
l'un sur les templates et l'autre sur `catalog`, sans conflits sur un énorme `main.go`. Les packages ne
dépendent que dans un sens (`server` utilise `catalog`, qui utilise `api`), donc je peux changer l'API
sans toucher au HTML. Enfin, charger les données une seule fois rend les pages plus rapides et ne
surcharge pas l'API : avant, chaque visite lui envoyait une ou deux requêtes.
