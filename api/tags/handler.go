// Package tags posts and deletes a tag, also gets all tags.
package tags

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler posts a tag and returns all the tags from the begining of the kind.
// Also, deletes a tag by given id as encoded key
// and gets all the tags from the begining of the kind.
func Handler(s *session.Session) {
	switch s.R.Method {
	case "POST":
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
		rb := new(api.ResponseBody)
		rb.Result = ts
		rb.PrevPageURL = "/tags?c=" + crsr.String()
		s.W.Header().Set("Content-Type", "application/json")
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
		ts, err := tag.GetMulti(s, nil)
		if err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = ts
		api.WriteResponseJSON(s, rb)
	}
}
