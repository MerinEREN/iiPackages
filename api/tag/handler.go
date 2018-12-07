// Package tag deletes a tag by the given id as encoded datastore key.
package tag

import (
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
	"strings"
)

// Handler handles delete request to delete a tag by the given id as encoded datastore key.
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	switch s.R.Method {
	case "DELETE":
		if ID == "" {
			log.Printf("Path: %s, Error: no tag ID to delete\n", s.R.URL.Path)
			http.Error(s.W, "No tag ID to delete", http.StatusBadRequest)
			return
		}
		err := tag.Delete(s.Ctx, ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	}
}
