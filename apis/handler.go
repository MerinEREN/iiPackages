/*
Package apis "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package apis

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
)

// WriteResponse "Exported functions should have a comment"
func WriteResponse(s *session.Session, rb *ResponseBody) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	s.W.Write(bs)
}
