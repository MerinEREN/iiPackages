// Package role deletes a role by the given id as encoded datastore key.
package role

import (
	"github.com/MerinEREN/iiPackages/datastore/role"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
	"strings"
)

// Handler handles delete request to delete a role by the given id as encoded datastore key.
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	if ID == "" {
		log.Printf("Path: %s, Error: no role ID\n", s.R.URL.Path)
		http.Error(s.W, "No role ID", http.StatusBadRequest)
		return
	}
	switch s.R.Method {
	case "DELETE":
		err := role.Delete(s.Ctx, ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "PUT" requests
	}
}
