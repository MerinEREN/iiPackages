// Package demand updates and returns a demand.
package demand

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/*
Handler patches status as "deleted" if method is "DELETE", updates and returns demand
via demand ID which is an encoded key.
*/
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	if ID == "" {
		log.Printf("Path: %s, Error: no demand ID\n", s.R.URL.Path)
		http.Error(s.W, "No demand ID", http.StatusBadRequest)
		return
	}
	switch s.R.Method {
	case "PUT":
		ct := s.R.Header.Get("Content-Type")
		if ct != "application/json" {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
				"Content type is not application/json")
			http.Error(s.W, "Content type is not application/json",
				http.StatusUnsupportedMediaType)
			return
		}
		bs, err := ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		d := new(demand.Demand)
		err = json.Unmarshal(bs, d)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		d, err = demand.Update(s.Ctx, d, ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.Header().Set("Content-Type", "application/json")
		api.WriteResponseJSON(s, d)
	case "DELETE":
		err := demand.UpdateStatus(s.Ctx, ID, "deleted")
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		k, err := datastore.DecodeKey(ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		d := new(demand.Demand)
		err = datastore.Get(s.Ctx, k, d)
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		d.ID = ID
		d.UserID = k.Parent().Encode()
		d.AccountID = k.Parent().Parent().Encode()
		s.W.Header().Set("Content-Type", "application/json")
		api.WriteResponseJSON(s, d)
	}
}
