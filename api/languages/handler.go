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
)

// Handler posts a language and returns limited languages from the begining of the kind.
// Also, deletes languages by given ids as encoded keys
// and gets languages from the begining of the given cursor..
// ADD AUTHORISATION CONTROL !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	switch s.R.Method {
	case "POST":
		langCode := s.R.FormValue("ID")
		contentID := s.R.FormValue("contentID")
		lang := &language.Language{
			ID:        langCode,
			ContentID: contentID,
		}
		mpf, hdr, err := s.R.FormFile("file")
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
		} else {
			defer mpf.Close()
			lang.Link, err = storage.UploadFile(s, mpf, hdr)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		var langs language.Languages
		langs, err = language.PutAndGetAll(s, lang)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
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
		lcx := q["IDs"]
		if len(lcx) == 0 {
			log.Printf("Path: %s, Error: no language code\n", URL.Path)
			http.Error(s.W, "no language code", http.StatusBadRequest)
			return
		} else if len(lcx) == 1 {
			k := datastore.NewKey(s.Ctx, "Language", lcx[0], 0, nil)
			err = datastore.Delete(s.Ctx, k)
		} else {
			kx := make([]*datastore.Key, len(lcx), len(lcx))
			for _, v := range lcx {
				kx = append(kx, datastore.NewKey(s.Ctx, "Language", v, 0, nil))
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
			return
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
