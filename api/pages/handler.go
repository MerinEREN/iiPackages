// Package pages posts a page, deletes and gets pages.
package pages

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/page"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handler posts a page and returns limited pages from the begining of the kind.
// Also, deletes pages by given ids as encoded keys
// and gets pages from the begining of the given cursor.
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
		p := new(page.Page)
		err = json.Unmarshal(bs, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		p.Created = time.Now()
		p.LastModified = time.Now()
		stringID := strings.ToLower(strings.Replace(p.Name, " ", "", -1))
		k := datastore.NewKey(s.Ctx, "Page", stringID, 0, nil)
		_, err = datastore.Put(s.Ctx, k, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		p.ID = k.Encode()
		sx := []string{"/pages", p.ID}
		path := strings.Join(sx, "/")
		rel, err := URL.Parse(path)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.Header().Set("Location", rel.String())
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponseJSON(s, p)
		// Using 'decoder' is an alternative and can be used if response body has
		// more than one json object.
		// Otherwise don't use it, because it has performance disadvantages
		// compared to first solution.
		/*decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		} */
		/*
			name := s.R.FormValue("name")
			p := &page.Page{
				Name:         name,
				Created:      time.Now(),
				LastModified: time.Now(),
			}
			mpf, hdr, err := s.R.FormFile("file")
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			} else {
				defer mpf.Close()
				p.Link, err = storage.UploadFile(s, mpf, hdr)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		*/
	case "DELETE":
		var err error
		q := URL.Query()
		ekx := q["IDs"]
		switch len(ekx) {
		case 0:
			log.Printf("Path: %s, Error: no page ids\n", URL.Path)
			http.Error(s.W, "no page ids", http.StatusBadRequest)
			return
		case 1:
			err = page.Delete(s.Ctx, ekx[0])
		default:
			err = page.DeleteMulti(s.Ctx, ekx)
		}
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		ps, err := page.GetAll(s.Ctx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		switch len(ps) {
		case 0:
			s.W.WriteHeader(http.StatusNoContent)
		case 1:
			s.W.Header().Set("Content-Type", "application/json")
			for _, v := range ps {
				api.WriteResponseJSON(s, v)
			}
		default:
			s.W.Header().Set("Content-Type", "application/json")
			api.WriteResponseJSON(s, ps)
		}
	}
}
