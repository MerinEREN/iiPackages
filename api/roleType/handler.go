// Package roleType deletes a roleType by the given id as datastore key's stringID.
package roleType

import (
	"github.com/MerinEREN/iiPackages/datastore/roleType"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
	"strings"
)

// Handler handles delete request to delete a roleType by the given id as encoded datastore key.
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	switch s.R.Method {
	case "DELETE":
		if ID == "" {
			log.Printf("Path: %s, Error: no roleType ID to delete\n", s.R.URL.Path)
			http.Error(s.W, "No roleType ID to delete", http.StatusBadRequest)
			return
		}
		err := roleType.Delete(s.Ctx, ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	}
}
