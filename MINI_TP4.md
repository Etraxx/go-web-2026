# Mini TP 4 : les vraies pages de GameDex

Objectif : GameDex affiche enfin de vraies pages HTML, construites avec des templates
(`html/template`), au lieu de simples lignes de texte.

## Ce que j'ai fait

- Un dossier `gamedex/templates/` avec `header.html` et `footer.html` (communs à toutes les pages,
  avec le lien vers mon CSS), puis `accueil.html`, `liste.html`, `jeu.html` et `erreur.html`.
  Chaque page commence par `{{template "header" .}}` et finit par `{{template "footer" .}}`.
- Les templates sont lus une seule fois, au démarrage :
  `var tmpl = template.Must(template.ParseGlob("templates/*.html"))`. Le chemin est relatif, donc
  GameDex se lance depuis le dossier `gamedex/`.
- Une struct par page : `ListePage` (Games, Query), `JeuPage` (Game, Tournaments) et
  `ErreurPage` (Code, Message).
- `/jeux` appelle `FetchGames()` et affiche les 12 jeux avec `{{range .Games}}`, chaque nom étant un
  lien vers `/jeux/{id}`.
- `/jeux/{id}` garde les vérifications du mini TP 2 (400 et 404), appelle `FetchGames()` et
  `FetchTournaments()`, retrouve le jeu avec une boucle `for` et garde les tournois où
  `t.GameID == id`.
- Une fonction `afficherErreur(w, code, message)` qui fait `w.WriteHeader(code)` **avant** d'écrire
  la page `erreur.html`. Une route attrape-tout `GET /` affiche la page 404 pour les adresses inconnues.
- J'ai adapté le mini TP 3 : plus d'appel à l'API ni de `log.Fatalf` au démarrage. Le site tourne
  même si l'API est éteinte, et chaque requête vers `/jeux` ou `/jeux/{id}` appelle l'API.
- `?q=` est lu avec `r.URL.Query().Get("q")` et affiché (« Vous avez cherché : … ») avec un petit
  formulaire GET. La liste n'est pas encore filtrée.

## Les codes obtenus

Mini API lancée (`go run ./api`), puis GameDex depuis `gamedex/` (`go run .`) :

| Requête | Code | Ce qui s'affiche |
| --- | --- | --- |
| `GET /` | 200 | Accueil avec header, footer et CSS |
| `GET /jeux` | 200 | Les 12 jeux (j'ai compté 12 liens `<a href="/jeux/N">`) |
| `GET /jeux/3` | 200 | Rocket League, 2015, Psyonix, sport et course, 3 tournois |
| `GET /jeux/2` | 200 | Hades : « Aucun tournoi prévu » |
| `GET /jeux/abc` | 400 | Page d'erreur |
| `GET /jeux/0` | 400 | Page d'erreur |
| `GET /jeux/-4` | 400 | Page d'erreur |
| `GET /jeux/999` | 404 | Page d'erreur |
| `GET /nimporte/quoi` | 404 | Page d'erreur (et pas 200) |
| `POST /jeux` | 405 | `Allow: GET, HEAD` |
| `GET /jeux?q=cel` | 200 | « Vous avez cherché : cel » |
| `GET /static/css/style.css` | 200 | `text/css; charset=utf-8` |
| `/jeux` et `/jeux/3`, API éteinte | 502 | Page d'erreur, en 0,001 s |
| `/jeux` et `/jeux/3`, avec `/api/slow` | 502 | Page d'erreur, en 3,00 s |

Pour `/jeux/3`, j'ai recoupé avec `curl http://localhost:8081/api/games/3/tournaments` :

```
Lyon Rocket Cup · Lyon · 2027-02-13
Aerial Masters · Madrid · 2027-04-03
Boost Open · Copenhague · 2027-06-19
```

Mêmes noms, mêmes villes et mêmes dates que dans le JSON. Pour trouver un jeu sans tournoi, j'ai
regardé les `gameId` de `/api/tournaments` : 1, 3, 4, 6, 7, 8, 9, 10 et 12. Les jeux 2, 5 et 11
n'ont donc pas de tournoi, et `/jeux/2` affiche bien « Aucun tournoi prévu » (le `{{else}}` du `range`).

## Le test `<script>`

J'ai appelé `/jeux?q=<script>alert("piraté")</script>` (URL encodée avec
`curl -G --data-urlencode`). Le HTML reçu :

```html
<input type="text" name="q" value="&lt;script&gt;alert(&#34;piraté&#34;)&lt;/script&gt;" placeholder="Nom d'un jeu">
<p class="recherche-texte">Vous avez cherché : &lt;script&gt;alert(&#34;piraté&#34;)&lt;/script&gt;</p>
```

Il n'y a aucun `<script>` actif : le navigateur affiche le texte au lieu de l'exécuter.

Pour comparer, j'ai remplacé brièvement `html/template` par `text/template` (sans commiter) :

```html
<p class="recherche-texte">Vous avez cherché : <script>alert("piraté")</script></p>
```

Cette fois le script est dans la page et le navigateur l'exécuterait : c'est une faille XSS. J'ai
remis `html/template` tout de suite.

## API éteinte et API lente

**API éteinte** : j'ai arrêté la mini API en laissant GameDex tourner. `/jeux` et `/jeux/3`
répondent 502 avec la page « Le service de données est indisponible, réessayez plus tard ». Le
serveur ne plante pas et l'accueil répond toujours 200. L'erreur technique n'apparaît que dans le
terminal :

```
2026/11/25 09:45:07 FetchGames : Get "http://localhost:8081/api/games": dial tcp [::1]:8081: connect: connection refused
```

J'ai relancé l'API : `/jeux` et `/jeux/3` sont repassés à 200 sans relancer GameDex.

**API lente** : j'ai remplacé temporairement `/api/games` par `/api/slow` (réponse au bout de 5 s).
`/jeux` répond 502 en **3,00 s** et `/jeux/3` en **3,00 s**, grâce au timeout de 3 secondes du client :

```
2026/11/25 09:54:54 FetchGames : Get "http://localhost:8081/api/slow": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

Puis j'ai remis `/api/games` (vérifié avec `git diff` avant de commiter).

`gofmt -l .` n'affiche rien, `go vet ./...` et `go build ./...` passent.

## Ce que j'en retiens

Avec des templates, le HTML est dans des fichiers `.html` séparés du code Go : c'est plus lisible
que du HTML dans des `Printf`, et le header et le footer ne sont écrits qu'une fois pour toutes les
pages. J'utilise `html/template` et pas `text/template` parce que `html/template` échappe
automatiquement ce qu'il insère (`<` devient `&lt;`), ce qui empêche un utilisateur d'injecter du
JavaScript dans ma page. `{{define "nom"}}` crée un morceau de page réutilisable, `{{template "nom" .}}`
l'insère, `{{range}}` répète un bloc pour chaque élément d'une liste, et `{{else}}` dans un `range`
s'affiche quand la liste est vide. Enfin, quand l'API tombe, je réponds 502 (Bad Gateway) : ce
n'est pas la faute de l'utilisateur (ce serait un 4xx), c'est le serveur dont je dépends qui ne
répond pas correctement.
