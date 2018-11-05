// Package languages posts a page, deletes and gets pages.
package languages

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/language"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// Handler posts a language and returns limited languages from the begining of the kind.
// Also, deletes languages by given ids as encoded keys
// and gets languages from the begining of the given cursor..
func Handler(s *session.Session) {
	switch s.R.Method {
	case "POST":
		langCode := s.R.FormValue("ID")
		name := s.R.FormValue("name")
		lang := &language.Language{
			ID:   langCode,
			Name: name,
		}
		mpf, hdr, err := s.R.FormFile("file")
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		} else {
			defer mpf.Close()
			lang.Link, err = storage.UploadFile(s, mpf, hdr)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		var langs language.Languages
		langs, crsr, err = language.PutAndGetMulti(s, crsr, lang)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = langs
		rb.PrevPageURL = "/languages?c=" + crsr.String()
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	case "DELETE":
		var err error
		langCodesAsString := s.R.FormValue("IDs")
		lcx := strings.Split(langCodesAsString, ",")
		if len(lcx) == 0 {
			log.Printf("Path: %s, Error: no language code\n", s.R.URL.Path)
			http.Error(s.W, "no language code", http.StatusBadRequest)
			return
		} else if len(lcx) == 1 {
			err = language.Delete(s, lcx[0])
		} else {
			err = language.DeleteMulti(s, lcx)
		}
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		crsr, err := datastore.DecodeCursor(s.R.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		langs, crsr, err := language.GetMulti(s, crsr, nil)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.PrevPageURL = "/languages?c=" + crsr.String()
		rb.Result = langs
		api.WriteResponseJSON(s, rb)
	}
}
