/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package language

import (
	api "github.com/MerinEREN/iiPackages/apis"
	"github.com/MerinEREN/iiPackages/datastore/language"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

func Handler(s *session.Session) {
	if s.R.Method == "POST" {
		langCode := s.R.FormValue("code")
		mpf, hdr, err := s.R.FormFile("file")
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		defer mpf.Close()
		lang := &language.Language{
			Code: langCode,
			Mpf:  mpf,
			Hdr:  hdr,
		}
		lang.Link, err = storage.UploadFile(s, lang.Mpf, lang.Hdr)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		lang, err = language.Put(s.Ctx, lang)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		// RETURN MediaLink HERE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		s.W.WriteHeader(201)
		return
	} else {
		rb := new(api.ResponseBody)
		if s.R.FormValue("action") == "getCount" {
			c, err := language.GetCount(s.Ctx)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			rb.Result = c
		} else {
			c, err := datastore.DecodeCursor(s.R.FormValue("c"))
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			langs, c, err := language.GetMulti(s.Ctx, c)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			rb.PrevPageURL = "/languages?d=prev&" + "c=" + c.String()
			rb.Result = langs
		}
		api.WriteResponse(s, rb)
	}
}
