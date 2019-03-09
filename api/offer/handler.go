// Package offer updates and returns a offer.
package offer

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/offer"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// Handler updates and returns offer via offer ID which is an encoded key.
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	if ID == "" {
		log.Printf("Path: %s, Error: no offer ID\n", s.R.URL.Path)
		http.Error(s.W, "No offer ID", http.StatusBadRequest)
		return
	}
	o := new(offer.Offer)
	k := new(datastore.Key)
	var err error
	switch s.R.Method {
	case "DELETE":
		err = offer.UpdateStatus(s.Ctx, ID, "deleted")
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
		return
	case "PUT":
		err = offer.UpdateStatus(s.Ctx, ID, "accepted")
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
		return
	default:
		k, err = datastore.DecodeKey(ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		err = datastore.Get(s.Ctx, k, o)
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
	}
	os := make(offer.Offers)
	o.ID = ID
	uKey, err := datastore.DecodeKey(o.UserID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	o.AccountID = uKey.Parent().Encode()
	os[ID] = o
	rb := new(api.ResponseBody)
	rb.Result = os
	api.WriteResponseJSON(s, rb)
}
