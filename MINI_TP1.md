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

## 4. Les tournois du jeu 4

```
curl.exe -i http://localhost:8081/api/games/4/tournaments
```

Le jeu 4 est Street Fighter 6. Il a **4 tournois** :

| id | Nom | Ville | Date |
| --- | --- | --- | --- |
| 4 | Fight Night Marseille | Marseille | 2027-01-30 |
| 5 | Combo Breaker EU | Paris | 2027-03-20 |
| 6 | Rising Dragon | Tokyo | 2027-08-07 |
| 7 | Hadoken Series | Bruxelles | 2027-05-15 |

Villes : Marseille, Paris, Tokyo et Bruxelles.

## 5. Filtres sur /api/games

```
curl.exe -i "http://localhost:8081/api/games?genre=course"
curl.exe -i "http://localhost:8081/api/games?year=2018"
curl.exe -i "http://localhost:8081/api/games?q=cel"
curl.exe -i "http://localhost:8081/api/games?q=CEL"
curl.exe -i "http://localhost:8081/api/games?genre=course&year=2020"
curl.exe -i "http://localhost:8081/api/games?q=cel&genre=course&year=2018"
```

| Filtre | Nombre de jeux | Jeux |
| --- | --- | --- |
| `genre=course` | 3 | Rocket League, Trackmania, Mario Kart 8 Deluxe |
| `year=2018` | 2 | Celeste, Dead Cells |
| `q=cel` | 2 | Celeste, Dead Cells (« cel » est aussi dans « Cells ») |
| `q=CEL` | 2 | Celeste, Dead Cells (pas de différence avec les majuscules) |
| `genre=course&year=2020` (cumulé) | 1 | Trackmania |
| `q=cel&genre=course&year=2018` (cumulé) | 0 | `[]` : une liste vide, mais quand même un code 200 |

Quand on cumule les filtres, un jeu doit respecter **tous** les filtres en même temps. Si rien ne
correspond, l'API ne renvoie pas d'erreur : elle renvoie une liste vide avec un 200.

## 6. /api/stats et la clé X-API-Key

```
curl.exe -i http://localhost:8081/api/stats
curl.exe -i -H "X-API-Key: eleve" http://localhost:8081/api/stats
curl.exe -i -H "X-API-Key: prof" http://localhost:8081/api/stats
```

| En-tête envoyé | Code | Réponse |
| --- | --- | --- |
| aucun | 401 Unauthorized | `{"error": "clé absente ou inconnue : ajoutez l'en-tête X-API-Key"}` |
| `X-API-Key: eleve` | 403 Forbidden | `{"error": "clé reconnue, mais pas le droit de voir les statistiques"}` |
| `X-API-Key: prof` | 200 OK | `{"games": 12, "tournaments": 20}` |

J'ai aussi essayé une clé inventée (`X-API-Key: toto`) : c'est aussi un 401.

Différence entre 401 et 403 : avec un 401, le serveur ne sait pas qui je suis (pas de clé ou clé
inconnue), alors qu'avec un 403 il sait qui je suis mais je n'ai pas le droit d'accéder à la ressource.

## 7. Bonus : /api/broken et /api/slow

```
curl.exe -i http://localhost:8081/api/broken
curl.exe -s -o NUL -w "code=%{http_code} time_total=%{time_total}s" http://localhost:8081/api/slow
Measure-Command { curl.exe -s -o NUL http://localhost:8081/api/slow }
```

- `/api/broken` répond toujours **500 Internal Server Error** avec `{"error": "panne simulée"}`.
  C'est une erreur du serveur, pas de ma requête.
- `/api/slow` répond **200 OK** avec la liste des jeux, mais seulement après environ 5 secondes.
  J'ai mesuré `time_total=5.000909s` avec curl.exe et `5,03 s` avec `Measure-Command`.
  Ça servira plus tard pour tester les timeouts d'un client HTTP.
