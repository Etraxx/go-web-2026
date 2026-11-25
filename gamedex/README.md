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

GameDex va chercher les jeux et les tournois dans la mini API : il faut donc d'abord lancer la mini
API, dans un premier terminal, depuis la racine du dépôt :

```
go run ./api
```

Puis, dans un second terminal, **depuis le dossier `gamedex/`** (les templates sont lus dans
`templates/`, un chemin relatif au dossier de lancement) :

```
cd gamedex
go run .
```

Le serveur écoute sur http://localhost:8080. Pour choisir un autre port :

```
go run . -port 9090
```

GameDex démarre même si la mini API est éteinte. Chaque visite de `/jeux` ou `/jeux/{id}` appelle
l'API (`/api/games` et `/api/tournaments`). Si elle est éteinte, trop lente (plus de 3 secondes),
répond autre chose que 200 ou renvoie un JSON illisible, la page d'erreur s'affiche avec le code
502. L'erreur technique n'est écrite que dans le terminal.

## Le dossier templates/

Les pages sont des templates `html/template`, chargés une seule fois au démarrage :

| Fichier | Rôle |
| --- | --- |
| `header.html` | `{{define "header"}}` : début du HTML, lien vers le CSS, en-tête avec les liens GameDex et Les jeux |
| `footer.html` | `{{define "footer"}}` : pied de page et fin du HTML |
| `accueil.html` | La page d'accueil |
| `liste.html` | La liste des jeux, le formulaire de recherche et « Vous avez cherché : … » |
| `jeu.html` | La fiche d'un jeu et ses tournois (ou « Aucun tournoi prévu ») |
| `erreur.html` | La page d'erreur (code et message) |

## Les routes

| Route | Ce qu'elle renvoie |
| --- | --- |
| `GET /` | La page d'accueil |
| `GET /jeux` | La liste des 12 jeux, chaque nom est un lien vers sa fiche |
| `GET /jeux?q=...` | La même liste, avec « Vous avez cherché : ... » au-dessus (pas encore de filtre) |
| `GET /jeux/{id}` | La fiche du jeu et ses tournois (nom, ville, date) |
| `GET /static/...` | Les fichiers du dossier `static/`, par exemple `/static/css/style.css` |
| Toute autre adresse en GET | La page d'erreur avec le code 404 |

Les codes d'erreur, tous affichés avec `erreur.html` :

| Cas | Code |
| --- | --- |
| Identifiant qui n'est pas un nombre ou plus petit que 1 (`/jeux/abc`, `/jeux/0`) | 400 |
| Jeu inexistant (`/jeux/999`) ou adresse inconnue | 404 |
| Mini API indisponible (éteinte, trop lente, erreur, JSON illisible) | 502 |

Une autre méthode que GET sur une route existante (par exemple `POST /jeux`) répond 405.
