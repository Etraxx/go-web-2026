package store

import "sync"

// Les inscriptions aux tournois, sur le même modèle que les likes :
// une map protégée par son propre mutex (numéro du tournoi -> pseudos inscrits).
var (
	muInscriptions sync.Mutex
	inscriptions   = map[int][]string{}
)

// Register inscrit le pseudo au tournoi. S'il est déjà inscrit, on ne l'ajoute pas une deuxième fois.
func Register(tournamentID int, pseudo string) {
	muInscriptions.Lock()
	defer muInscriptions.Unlock()

	// On regarde si le pseudo est déjà dans la liste
	inscrits := inscriptions[tournamentID]
	for i := 0; i < len(inscrits); i++ {
		if inscrits[i] == pseudo {
			return
		}
	}
	inscriptions[tournamentID] = append(inscrits, pseudo)
}

// Registrations renvoie une copie de la liste des inscrits d'un tournoi.
// C'est une copie : celui qui la reçoit ne peut pas modifier la liste gardée ici.
func Registrations(tournamentID int) []string {
	muInscriptions.Lock()
	defer muInscriptions.Unlock()

	inscrits := inscriptions[tournamentID]
	copie := make([]string, len(inscrits))
	copy(copie, inscrits)
	return copie
}
