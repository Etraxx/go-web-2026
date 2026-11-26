# Mini TP 6 : filtrer et liker

Objectif : filtrer les jeux avec un formulaire GET (c'est la mini API qui filtre), et liker un jeu
avec un formulaire POST, sans que F5 refasse le like.

Une première version de ce mini TP avait été faite à partir de l'ancien énoncé. Je l'ai annulée
(commits retirés), car l'énoncé du cours a changé.

## Ce que j'ai fait

- **`internal/api`** : `SearchGames(q, genre)` construit les paramètres avec `url.Values`
  (`params.Set("q", q)`, `params.Set("genre", genre)`) et appelle `/api/games?` + `params.Encode()`
  avec le client à 3 secondes. Comme `FetchGames`, elle vérifie l'erreur et le code 200, puis
  décode le JSON.
- **`internal/catalog`** : `Genres(games)` renvoie tous les genres sans doublon (boucle et map),
  triés avec `sort.Strings`. `catalog` n'importe que `gamedex/internal/api` et `sort`.
- **`internal/server`**, handler `listeJeux` :
  - il lit `q` et `genre` dans l'URL et appelle `api.SearchGames(q, genre)` : c'est la mini API qui
    filtre ;
  - si l'appel échoue, il fait un `log.Println` de l'erreur et affiche la page d'erreur 502 « Le
    service des jeux ne répond pas. » ;
  - la liste des genres est calculée sur `s.Games` (tous les jeux chargés au démarrage), pas sur le
    résultat filtré.
- **`liste.html`** : le formulaire GET `class="filtres"` (champ `q` et `<select name="genre">`, avec
  `{{if eq . $.Genre}}selected{{end}}`), « N jeu(x) trouvé(s) », la liste, et « Aucun jeu ne
  correspond. » quand elle est vide. Dans le `range`, `.` est le genre en cours et `$` la donnée de
  toute la page, d'où `$.Genre`.
- **`internal/store/likes.go`** : une map numéro du jeu → nombre de likes, protégée par un
  `sync.Mutex`, avec `AddLike(gameID)` et `Likes(gameID)`. Le fichier `snippets/likes.go` du cours
  n'existait pas dans le dépôt d'origine (une seule branche `master`, un seul commit « Init »), donc
  je l'ai écrit moi-même sur le modèle de l'énoncé.
- **Fiche d'un jeu** : un formulaire POST `/jeux/{{.Game.ID}}/like` avec le bouton
  « ♥ J'aime (N) ». `JeuPage` a un champ `Likes` rempli avec `store.Likes(id)`.
- **Route `POST /jeux/{id}/like`** : elle vérifie l'id comme la fiche (400 ou 404, grâce à une petite
  fonction `chercherJeu` commune aux deux handlers), appelle `store.AddLike`, puis redirige avec
  `http.Redirect(..., http.StatusSeeOther)`, c'est-à-dire un 303 vers `/jeux/{id}`.

## Les filtres

Mini API lancée avec sa sortie dans un fichier log, puis GameDex. Les accents sont envoyés avec
`curl -G --data-urlencode`.

| URL GameDex | Jeux obtenus | Mini API directement (port 8081) |
| --- | --- | --- |
| `/jeux` | 12 jeu(x) trouvé(s) : les 12 jeux | |
| `/jeux?q=cel` | Celeste, Dead Cells | `?q=cel` : Celeste, Dead Cells |
| `/jeux?q=CEL` | Celeste, Dead Cells | |
| `/jeux?genre=course` | Rocket League, Trackmania, Mario Kart 8 Deluxe | mêmes 3 jeux |
| `/jeux?genre=indé&q=ce` | Celeste, Dead Cells | mêmes 2 jeux |
| `/jeux?genre=zzz` | 200, 0 jeu(x), « Aucun jeu ne correspond. » | |

Je m'attendais à ce que `q=cel` ne donne que Celeste. Mais « Dead **Cel**ls » contient aussi « cel »,
et la mini API renvoie elle-même les deux jeux : le résultat est donc juste.

Le `<select>` contient les 16 genres sans doublon et triés (action, combat, course, fps, indé,
metroidvania, moba, party, plateforme, rogue-like, sandbox, simulation, sport, stratégie, survie,
tactique). Il garde le genre choisi même quand le résultat est filtré, par exemple avec
`/jeux?genre=indé&q=ce` :

```html
<option value="indé" selected>indé</option>
```

Le log de la mini API montre les recherches transmises par GameDex. Les 2 premières lignes viennent du
démarrage, les autres de chaque visite de `/jeux` :

```
2026/11/26 09:38:40 GET /api/games
2026/11/26 09:38:40 GET /api/tournaments
2026/11/26 09:38:41 GET /api/games?genre=&q=
2026/11/26 09:38:41 GET /api/games?genre=&q=cel
2026/11/26 09:38:41 GET /api/games?genre=&q=CEL
2026/11/26 09:38:41 GET /api/games?genre=course&q=
2026/11/26 09:38:41 GET /api/games?genre=ind%C3%A9&q=ce
2026/11/26 09:38:41 GET /api/games?genre=zzz&q=
```

`url.Values` encode le « é » en `%C3%A9` et met les paramètres dans l'ordre alphabétique (`genre`
puis `q`). Les fiches `/jeux/{id}` n'apparaissent pas dans ce log : elles utilisent les données en
mémoire.

L'échappement marche toujours. Avec `/jeux?q=<script>alert(1)</script>` (encodé) :

```html
<p class="recherche-texte">Vous avez cherché : &lt;script&gt;alert(1)&lt;/script&gt;</p>
```

**API éteinte** (GameDex déjà démarré) : `/jeux?q=a` répond `502 Bad Gateway` avec la page d'erreur
« Le service des jeux ne répond pas. ». L'erreur technique n'apparaît que dans le terminal :

```
2026/11/26 09:38:50 Get "http://localhost:8081/api/games?genre=&q=a": dial tcp [::1]:8081: connect: connection refused
```

`/jeux/3` répond toujours 200 (données en mémoire) et GameDex ne plante pas. Une fois l'API
rallumée, `/jeux` répond de nouveau 200.

## Les likes : POST → 303 → GET, et le F5

```
GET /jeux/3                     ♥ J'aime (0)
$ curl -i -X POST http://localhost:8080/jeux/3/like
HTTP/1.1 303 See Other
Location: /jeux/3
GET /jeux/3                     ♥ J'aime (1)
GET /jeux/3 (le « F5 »)         ♥ J'aime (1)    <- le compteur ne remonte pas
2e POST, puis GET               ♥ J'aime (2)
après 10 POST en tout           ♥ J'aime (10)
GET /jeux/4                     ♥ J'aime (0)    <- les likes du jeu 3 ne touchent pas le jeu 4
```

Pour imiter le navigateur, j'ai aussi lancé `curl -L -d "" .../jeux/4/like` : curl suit le 303 en
faisant un GET et termine sur `/jeux/4` avec un code 200. (Avec `curl -L -X POST`, curl garde la
méthode POST après la redirection et tombe sur 405 : c'est parce que `-X` force la méthode, un
navigateur ne fait pas ça.)

Erreurs :

| Requête | Code |
| --- | --- |
| `POST /jeux/abc/like` | 400 |
| `POST /jeux/0/like` | 400 |
| `POST /jeux/999/like` | 404 |
| `GET /jeux/3/like` | 404 |

`GET /jeux/3/like` répond 404 et pas 405. La route attrape-tout `GET /` accepte toutes les adresses en
GET, y compris `/jeux/3/like`, et elle affiche ma page 404. Comme une route GET correspond, Go ne
répond pas « méthode non autorisée ». Je l'ai laissé comme ça.

**Concurrence** : j'ai lancé GameDex compilé avec `-race` (le détecteur d'accès concurrents de Go),
puis envoyé 20 POST en parallèle sur le jeu 5. Résultat : exactement `♥ J'aime (20)`. J'ai ensuite
envoyé 50 POST et 50 GET mélangés, en parallèle : `♥ J'aime (70)`. Le détecteur n'a signalé aucun
« DATA RACE ».

Les anciennes pages marchent toujours : `/` 200, `/jeux/abc` 400, `/jeux/999` 404, `/nimporte/quoi`
404, `POST /jeux` 405 et `/static/css/style.css` 200. `gofmt -l .` n'affiche rien, `go vet ./...` et
`go build ./...` passent, et `catalog` n'importe ni `net/http`, ni `html/template`, ni `store`.

## Ce que j'en retiens

On utilise GET pour **lire** (filtrer, chercher) : les champs vont dans l'URL, on peut la partager ou
la recharger sans risque. On utilise POST pour **agir** (liker) : les champs vont dans le corps de la
requête. Après le POST, je redirige avec un 303 vers la fiche : le navigateur fait alors un GET, et
c'est ce GET qu'un F5 recharge, donc le like n'est pas refait (« post / redirect / get »). Le mutex
sert parce que le serveur traite plusieurs requêtes en même temps : sans lui, deux likes simultanés
pourraient modifier la map en même temps et un like serait perdu, ou le programme planterait. Les
likes sont dans une map en mémoire, donc ils disparaissent quand le serveur s'arrête. Une vraie base
de données (SQLite, PostgreSQL…) les écrirait sur le disque pour qu'ils survivent à un redémarrage.

## Ce que je n'ai pas pu vérifier

Je n'ai pas cliqué moi-même sur le bouton J'aime dans un vrai navigateur, ni appuyé sur F5. J'ai imité
le navigateur avec curl (POST, puis le GET de la redirection), et j'ai seulement vu l'affichage du
bouton et du formulaire dans des captures de Chrome en mode headless. Il faut encore faire le test à
la main : cliquer sur J'aime, puis F5.

## Bonus 1 : les plus likés en premier

- Dans `internal/catalog`, `SortByLikes(games, likes)` trie les jeux par nombre de likes, du plus
  grand au plus petit, avec `sort.SliceStable`. « Stable » veut dire qu'à égalité, l'ordre de l'API
  est gardé.
- `catalog` ne connaît pas `store` : c'est le handler `listeJeux` qui construit la map
  (`likes[g.ID] = store.Likes(g.ID)`, avec une boucle), puis qui appelle `catalog.SortByLikes`.
- Chaque jeu de la liste affiche « ♥ N » à côté de son nom (`{{index $.Likes .ID}}`).

Test : j'ai donné 3 likes à Trackmania, 1 à Hollow Knight et 1 à Valorant.

```
/jeux      : Trackmania ♥3, Valorant ♥1, Hollow Knight ♥1, Celeste ♥0, Hades ♥0, Rocket League ♥0, ...
/jeux?q=a  : Trackmania ♥3, Valorant ♥1, Hades ♥0, Rocket League ♥0, Stardew Valley ♥0, ...
```

Trackmania passe en premier. Valorant et Hollow Knight ont chacun 1 like : Valorant reste devant, comme
dans l'ordre de l'API. Les jeux à 0 gardent aussi l'ordre de l'API. Avec une recherche, le tri
s'applique au résultat renvoyé par la mini API. Les likes ne sont pas persistants : après un
redémarrage, tout revient à ♥0.

## Bonus 2 : inscription à un tournoi

- `internal/store/inscriptions.go` suit le même modèle que les likes : une map
  numéro du tournoi → pseudos, protégée par son propre mutex.
  - `Register` n'ajoute pas un pseudo déjà inscrit.
  - `Registrations` renvoie une **copie** de la liste (`make` + `copy`), pour que personne ne puisse
    modifier la liste interne sans passer par le mutex.
- `catalog.FindTournament` retrouve un tournoi par son numéro.
- Route `POST /tournois/{id}/inscription` :
  - id invalide : 400 ; tournoi inexistant : 404 ;
  - pseudo (`strings.TrimSpace`) vide ou de plus de 30 caractères : 400, avec le message « Le pseudo
    doit faire entre 1 et 30 caractères. » ;
  - sinon : `store.Register`, puis redirection 303 vers la fiche du jeu du tournoi.
- Sur la fiche, chaque tournoi affiche « N inscrit(s) · pseudo · pseudo », avec un petit formulaire
  POST (champ `pseudo` et bouton « S'inscrire »). `JeuPage.Tournaments` est maintenant une slice de
  `TournoiPage { Tournament ; Inscrits }`.

Tests sur le tournoi 1 (Lyon Rocket Cup, jeu 3) :

```
$ curl -i --data-urlencode "pseudo=Etraxx" http://localhost:8080/tournois/1/inscription
HTTP/1.1 303 See Other
Location: /jeux/3
fiche : 1 inscrit(s) · Etraxx
```

| Test | Résultat |
| --- | --- |
| même pseudo une 2e fois | toujours 1 inscrit |
| `"  Etraxx  "` (espaces autour) | toujours 1 inscrit (les espaces sont enlevés) |
| pseudo vide ou seulement des espaces | 400 |
| pseudo de 31 caractères | 400 (30 caractères : 303) |
| tournoi 9999 | 404 |
| tournoi `abc` | 400 |
| pseudo `<script>alert(1)</script>` | affiché `&lt;script&gt;alert(1)&lt;/script&gt;`, donc échappé par html/template |

Un jeu sans tournoi (Hades) n'affiche aucun formulaire d'inscription, seulement « Aucun tournoi prévu ».

Vérification finale : `gofmt -l .` n'affiche rien, `go vet ./...` et `go build ./...` passent, et
`catalog` n'importe ni `net/http`, ni `html/template`, ni `store`. Dans le log de la mini API, seules
les visites de `/jeux` arrivent : la fiche, le like et l'inscription n'appellent pas l'API.

J'ai vu l'affichage des inscrits et des formulaires dans une capture de Chrome en mode headless, mais
je n'ai pas rempli ni envoyé le formulaire d'inscription moi-même dans un vrai navigateur. Je l'ai
seulement fait avec curl.
