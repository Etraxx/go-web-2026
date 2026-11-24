# GameDex · votre travail

C'est ici que vous écrivez GameDex, mini TP après mini TP.

Au mini TP 2, placez-vous dans ce dossier et créez votre module :

```
cd gamedex
go mod init gamedex
```

Vous lancerez ensuite votre serveur depuis ce dossier (`go run .`, puis `go run ./cmd/gamedex` après le mini TP 5).

À la fin de chaque mini TP : commit, puis push sur votre fork.

## Lancer GameDex

GameDex interroge la mini API au démarrage : il faut donc d'abord lancer la mini API, dans un
premier terminal, depuis la racine du dépôt :

```
go run ./api
```

Puis, dans un second terminal, depuis la racine du dépôt :

```
cd gamedex
go run .
```

Au démarrage, GameDex récupère les jeux (`/api/games`) et les tournois (`/api/tournaments`) et
affiche dans le terminal chaque jeu avec son année et son nombre de tournois. Si la mini API est
éteinte, trop lente (plus de 3 secondes) ou répond autre chose que 200, GameDex affiche un message
d'erreur et s'arrête.

Le serveur écoute sur http://localhost:8080. Pour choisir un autre port :

```
go run . -port 9090
```

Les routes :

| Route | Ce qu'elle renvoie |
| --- | --- |
| `GET /` | La page d'accueil en HTML (avec le CSS) |
| `GET /jeux` | La liste des jeux (texte) |
| `GET /jeux/{id}` | La fiche d'un jeu (texte) : 400 si l'identifiant n'est pas un nombre ou est plus petit que 1, 404 s'il est plus grand que 12 |
| `GET /static/...` | Les fichiers du dossier `static/`, par exemple `/static/css/style.css` |

Toute autre adresse répond 404.
