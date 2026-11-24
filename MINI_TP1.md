# Mini TP 1 : explorer la mini API

Objectif : apprendre à lire une API avant de l'appeler depuis Go.

J'ai fait les mêmes requêtes que le bouton « Try it out » de Swagger, mais en ligne de commande avec
`curl.exe -i` (sous PowerShell il faut bien taper `curl.exe` et pas `curl`). L'option `-i` affiche le
code de réponse et les en-têtes avant le JSON, comme Swagger.

## 1. Lancer la mini API

Depuis la racine du dépôt :

```
go run ./api
```

Puis je vérifie la page Swagger :

```
curl.exe -s -o NUL -w "%{http_code} %{content_type}" http://localhost:8081/api/docs
```

Résultat : `200 text/html; charset=utf-8`. La documentation Swagger est bien en ligne.

## 2. Tous les endpoints

Toutes les réponses JSON ont le Content-Type `application/json; charset=utf-8`.

| Méthode | URL | Code | Extrait du JSON reçu |
| --- | --- | --- | --- |
| GET | `/api/games` | 200 OK | une liste de 12 jeux : `{"id": 1, "name": "Celeste", "year": 2018, "studio": "Maddy Makes Games", "genres": ["plateforme", "indé"]}`, ... |
| GET | `/api/games?genre=course` | 200 OK | Rocket League, Trackmania, Mario Kart 8 Deluxe |
| GET | `/api/games/3` | 200 OK | `{"id": 3, "name": "Rocket League", "year": 2015, "studio": "Psyonix", "genres": ["sport", "course"]}` |
| GET | `/api/games/4/tournaments` | 200 OK | une liste de 4 tournois : `{"id": 4, "gameId": 4, "name": "Fight Night Marseille", "city": "Marseille", "date": "2027-01-30", ...}`, ... |
| GET | `/api/tournaments` | 200 OK | une liste de 20 tournois : `{"id": 1, "gameId": 3, "name": "Lyon Rocket Cup", "city": "Lyon", ...}`, ... |
| GET | `/api/stats` (en-tête `X-API-Key: prof`) | 200 OK | `{"games": 12, "tournaments": 20}` |
| GET | `/api/broken` | 500 Internal Server Error | `{"error": "panne simulée"}` |
| GET | `/api/slow` | 200 OK (après 5 s) | la même liste de 12 jeux que `/api/games` |

Commandes utilisées :

```
curl.exe -i http://localhost:8081/api/games
curl.exe -i "http://localhost:8081/api/games?genre=course"
curl.exe -i http://localhost:8081/api/games/3
curl.exe -i http://localhost:8081/api/games/4/tournaments
curl.exe -i http://localhost:8081/api/tournaments
curl.exe -i -H "X-API-Key: prof" http://localhost:8081/api/stats
curl.exe -i http://localhost:8081/api/broken
curl.exe -i http://localhost:8081/api/slow
```

Remarque : dans PowerShell, je dois mettre l'URL entre guillemets quand il y a un `&` dans les filtres.

J'ai aussi essayé `curl.exe -i -X POST http://localhost:8081/api/games` : la réponse est
`405 Method Not Allowed` avec l'en-tête `Allow: GET, HEAD`, parce que la route n'accepte que GET.

## 3. Jeu qui n'existe pas et identifiant qui n'est pas un nombre

```
curl.exe -i http://localhost:8081/api/games/999
curl.exe -i http://localhost:8081/api/games/abc
```

| URL | Code | Message |
| --- | --- | --- |
| `/api/games/999` | 404 Not Found | `{"error": "jeu introuvable"}` |
| `/api/games/abc` | 400 Bad Request | `{"error": "l'identifiant doit être un nombre"}` |

999 est un nombre correct mais aucun jeu n'a cet id, donc 404 (la ressource n'existe pas).
`abc` n'est pas un nombre, c'est ma requête qui est mal formée, donc 400.
