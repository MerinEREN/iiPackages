/*
Package content "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package content

import (
	"github.com/MerinEREN/iiPackages/datastore/content"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
	"strings"
)

// Handler handles PUT requests for a content.
func Handler(s *session.Session) {
	var err error
	IDsAsString := s.R.FormValue("IDs")
	IDx := strings.Split(IDsAsString, ",")
	if len(IDx) == 1 {
		err = content.Delete(s, IDx[0])
	} else {
		err = content.DeleteMulti(s, IDx)
	}
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	s.W.WriteHeader(http.StatusNoContent)
}
