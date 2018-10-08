/*
Package tags "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package tags

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler "Exported functions should have a comment"
func Handler(s *session.Session) {
	switch s.R.Method {
	case "POST":
		rb := new(api.ResponseBody)
		name := s.R.FormValue("name")
		t := &tag.Tag{
			Name: name,
		}
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		ts, err := tag.PutAndGetMulti(s, t)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.Result = ts
		rb.Reset = true
		rb.PrevPageURL = "/tags?c=" + crsr.String()
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	case "DELETE":
		ID := s.R.FormValue("ID")
		if ID == "" {
			log.Printf("Path: %s, Error: no tag ID to delete\n", s.R.URL.Path)
			http.Error(s.W, "No tag ID to delete", http.StatusBadRequest)
			return
		}
		err := tag.Delete(s, ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		rb := new(api.ResponseBody)
		ts, err := tag.GetMulti(s, nil)
		if err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.Result = ts
		api.WriteResponse(s, rb)
	}
}
