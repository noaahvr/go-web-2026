package store

import (
	"errors"
	"slices"
	"strings"
	"sync"
)

var ErrDejaInscrit = errors.New("pseudo déjà inscrit à ce tournoi")

var (
	inscriptionsMu sync.RWMutex
	inscriptions   = make(map[int][]string)
)

func AddRegistration(tournamentID int, pseudo string) error {
	inscriptionsMu.Lock()
	defer inscriptionsMu.Unlock()
	for _, p := range inscriptions[tournamentID] {
		if strings.EqualFold(p, pseudo) {
			return ErrDejaInscrit
		}
	}
	inscriptions[tournamentID] = append(inscriptions[tournamentID], pseudo)
	return nil
}

func Registrations(tournamentID int) []string {
	inscriptionsMu.RLock()
	defer inscriptionsMu.RUnlock()
	return slices.Clone(inscriptions[tournamentID])
}
