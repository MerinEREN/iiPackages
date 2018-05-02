/*
Package languages "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
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

// Handler "Exported functions should have a comment"
func Handler(s *session.Session) {
	rb := new(api.ResponseBody)
	switch s.R.Method {
	case "POST":
		langCode := s.R.FormValue("ID")
		lang := &language.Language{
			ID: langCode,
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
		var c datastore.Cursor
		var langs language.Languages
		langs, c, err = language.PutAndGetMulti(s, c, lang)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.Result = langs
		rb.Reset = true
		rb.PrevPageURL = "/languages?c=" + c.String()
		s.W.WriteHeader(http.StatusCreated)
	case "DELETE":
		var err error
		langCodesAsString := s.R.FormValue("IDs")
		lcx := strings.Split(langCodesAsString, ",")
		if len(lcx) == 0 {
			log.Printf("Path: %s, Error: no language code\n", s.R.URL.Path)
			http.Error(s.W, "no language code", http.StatusBadRequest)
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
		if s.R.FormValue("action") == "getCount" {
			cnt, err := language.GetCount(s)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			rb.Result = cnt
		} else {
			c, err := datastore.DecodeCursor(s.R.FormValue("c"))
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			langs, c, err := language.GetMulti(s, c, nil)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			rb.PrevPageURL = "/languages?c=" + c.String()
			rb.Result = langs
		}
	}
	api.WriteResponse(s, rb)
}
