# Mini TP 3 : GameDex appelle la mini API

Objectif : GameDex ne se contente plus d'écrire du texte, il va chercher les jeux et les tournois
dans la mini API, au démarrage.

## Ce que j'ai ajouté dans gamedex/main.go

- Les structs `Game` et `Tournament`, avec des tags json qui correspondent exactement aux clés de
  l'API (`id`, `name`, `year`, `studio`, `genres` et `id`, `gameId`, `name`, `city`, `date`, `lat`, `lng`).
  J'ai vérifié les clés avec `curl.exe http://localhost:8081/api/games/1` et `/api/tournaments`.
- Un client HTTP avec un timeout de 3 secondes : `var client = &http.Client{Timeout: 3 * time.Second}`,
  et l'adresse de l'API dans une constante `adresseAPI`.
- `FetchGames` et `FetchTournaments`, qui passent par une petite fonction commune `lireJSON` :
  `client.Get`, vérification de `err`, `defer resp.Body.Close()`, vérification du code 200, puis
  `json.NewDecoder(resp.Body).Decode(...)`. Elles n'affichent rien, elles renvoient une erreur.
- Une map `byGame` (clé = `GameID`) pour compter les tournois de chaque jeu.
- Dans `main`, si un des deux Fetch échoue, j'affiche un message en français avec `log.Fatalf` et
  le programme s'arrête (pas de panic).

## 1. API allumée

Dans un premier terminal, à la racine : `go run ./api`. Dans un second : `cd gamedex` puis `go run .`

```
Jeux chargés depuis la mini API :
  Celeste (2018)               1 tournoi(s)
  Hades (2020)                 0 tournoi(s)
  Rocket League (2015)         3 tournoi(s)
  Street Fighter 6 (2023)      4 tournoi(s)
  Stardew Valley (2016)        0 tournoi(s)
  League of Legends (2009)     3 tournoi(s)
  Valorant (2020)              3 tournoi(s)
  Dead Cells (2018)            1 tournoi(s)
  Trackmania (2020)            2 tournoi(s)
  Mario Kart 8 Deluxe (2017)   2 tournoi(s)
  Hollow Knight (2017)         0 tournoi(s)
  Minecraft (2011)             1 tournoi(s)
Serveur sur http://localhost:8080
```

J'ai recoupé : 1 + 0 + 3 + 4 + 0 + 3 + 3 + 1 + 2 + 2 + 0 + 1 = 20, et
`curl.exe http://localhost:8081/api/tournaments` renvoie bien 20 tournois. Ensuite le serveur du mini
TP 2 marche comme avant : `/`, `/jeux` et `/jeux/3` répondent 200, `/jeux/abc` 400, `/jeux/999` 404,
le CSS 200, et `go run . -port 9090` marche aussi.

## 2. API coupée

J'ai arrêté la mini API puis lancé `go run .` :

```
2026/11/24 12:34:44 Impossible de charger les jeux : Get "http://localhost:8081/api/games": dial tcp [::1]:8081: connectex: No connection could be made because the target machine actively refused it.
exit status 1
```

Le programme s'arrête avec le code de sortie 1, sans panic ni stacktrace. Sous Windows, le
« connection refused » s'écrit « the target machine actively refused it » : personne n'écoute sur
le port 8081.

## 3. Avec /api/broken

J'ai remplacé temporairement `/api/games` par `/api/broken` dans `FetchGames` :

```
2026/11/24 12:34:00 Impossible de charger les jeux : l'API a répondu 500
exit status 1
```

Ici `client.Get` a renvoyé `err == nil` : pour Go, la requête a marché, puisque le serveur a bien
répondu. C'est mon test `resp.StatusCode != http.StatusOK` qui détecte le problème.

## 4. Avec /api/slow

J'ai remplacé temporairement `/api/games` par `/api/slow` (l'API répond après 5 secondes) :

```
2026/11/24 12:34:16 Impossible de charger les jeux : Get "http://localhost:8081/api/slow": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

Mesuré avec `Measure-Command` sur le programme compilé : il abandonne au bout de **3,05 s**, donc
juste après le timeout de 3 secondes, sans attendre les 5 secondes de l'API.

J'ai aussi essayé avec un timeout de 6 secondes : cette fois les 12 jeux s'affichent et le serveur
démarre au bout de 5,28 s. Puis j'ai remis `/api/games` et 3 secondes, et j'ai vérifié que tout
remarchait.

## Ce que j'en retiens

`http.Get` utilise le client par défaut, qui n'a pas de timeout : si l'API ne répond jamais, mon
programme attend pour toujours. Avec `http.Client{Timeout: 3 * time.Second}`, j'abandonne au bout de
3 secondes avec une erreur que je peux afficher. Ensuite, `err == nil` veut seulement dire que j'ai
reçu une réponse, pas que c'est une bonne réponse : une erreur 500 ou 404 arrive aussi avec
`err == nil`. Il faut donc toujours vérifier `resp.StatusCode` avant de lire le JSON.
