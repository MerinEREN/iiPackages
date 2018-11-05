// Package page gets, deletes and modifies a page.
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

// Handler gets, deletes and modifies the page by given page ID as encoded key.
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
		api.WriteResponseJSON(s, rb)
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
		api.WriteResponseJSON(s, rb)
	}
}
