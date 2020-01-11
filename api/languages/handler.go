// Package languages posts a page, deletes and gets pages.
package languages

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/language"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Handler posts a language and returns limited languages from the begining of the kind.
// Also, deletes languages by given ids as encoded keys
// and gets languages from the begining of the given cursor..
// ADD AUTHORISATION CONTROL !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Handler(s *session.Session) {
	URL := s.R.URL
	switch s.R.Method {
	case "POST":
		ct := s.R.Header.Get("Content-Type")
		if ct != "application/json" {
			log.Printf("Path: %s, Error: %v\n", URL.Path,
				"Content type is not application/json")
			http.Error(s.W, "Content type is not application/json",
				http.StatusUnsupportedMediaType)
			return
		}
		bs, err := ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		l := new(language.Language)
		err = json.Unmarshal(bs, l)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		l.Created = time.Now()
		k := datastore.NewKey(s.Ctx, "Language", l.ID, 0, nil)
		_, err = datastore.Put(s.Ctx, k, l)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponseJSON(s, l)
	case "DELETE":
		var err error
		q := URL.Query()
		lcx := q["IDs"]
		switch len(lcx) {
		case 0:
			log.Printf("Path: %s, Error: no language code\n", URL.Path)
			http.Error(s.W, "no language code", http.StatusBadRequest)
			return
		case 1:
			k := datastore.NewKey(s.Ctx, "Language", lcx[0], 0, nil)
			err = datastore.Delete(s.Ctx, k)
		default:
			kx := make([]*datastore.Key, len(lcx), len(lcx))
			for _, v := range lcx {
				kx = append(kx, datastore.NewKey(s.Ctx, "Language", v, 0,
					nil))
			}
			err = datastore.DeleteMulti(s.Ctx, kx)
		}
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		ls, err := language.GetAll(s.Ctx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		switch len(ls) {
		case 0:
			s.W.WriteHeader(http.StatusNoContent)
		case 1:
			s.W.Header().Set("Content-Type", "application/json")
			for _, v := range ls {
				api.WriteResponseJSON(s, v)
			}
		default:
			s.W.Header().Set("Content-Type", "application/json")
			api.WriteResponseJSON(s, ls)
		}
	}
}
