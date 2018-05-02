/*
Package page "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package page

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/page"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// Handler "Exported functions should have a comment"
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	switch s.R.Method {
	case "PUT":
		// GET THE ENTITY AS BLOB OBJECT TO PREVENT RunInTransaction IN PUT FUNC !!
		title := s.R.FormValue("title")
		p := &page.Page{
			ID:    ID,
			Title: title,
		}
		mpf, hdr, err := s.R.FormFile("file")
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		} else {
			defer mpf.Close()
			// CHECK THE STORAGE AND IF THE FILE PRESENT DO NOT UPLOAD THE FILE
			p.Link, err = storage.UploadFile(s, mpf, hdr)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		p, err = page.Put(s, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		ps := make(page.Pages)
		ps[p.ID] = p
		rb.Result = ps
		// s.W.WriteHeader(http.StatusOK)
		api.WriteResponse(s, rb)
	case "DELETE":
		err := page.Delete(s, ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		// REDIRECT TO THE PAGES !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		ps, err := page.Get(s, ID)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = ps
		api.WriteResponse(s, rb)
	}
}
