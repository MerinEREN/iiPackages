/*
Package language "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package language

import (
	"github.com/MerinEREN/iiPackages/datastore/language"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
)

// Handler "Exported functions should have a comment"
func Handler(s *session.Session) {
	switch s.R.Method {
	case "PUT":
	case "DELETE":
		langCode := s.R.FormValue("code")
		err := language.Delete(s, langCode)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
	}
}
