// Package pages posts a page, deletes and gets pages.
package pages

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

// Handler posts a page and returns limited pages from the begining of the kind.
// Also, deletes pages by given ids as encoded keys
// and gets pages from the begining of the given cursor.
func Handler(s *session.Session) {
	switch s.R.Method {
	case "POST":
		title := s.R.FormValue("title")
		p := &page.Page{
			Title: title,
		}
		mpf, hdr, err := s.R.FormFile("file")
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		} else {
			defer mpf.Close()
			p.Link, err = storage.UploadFile(s, mpf, hdr)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		/* bs, err := ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		p := new(page.Page)
		err = json.Unmarshal(bs, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		} */
		// Using 'decoder' is an alternative and can be used if response body has
		// more than one json object.
		// Otherwise don't use it, because it has performance disadvantages
		// compared to first solution.
		/*decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		} */
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		var ps page.Pages
		ps, crsr, err = page.PutAndGetMulti(s, crsr, p)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = ps
		rb.PrevPageURL = "/pages?c=" + crsr.String()
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	case "DELETE":
		IDsAsString := s.R.FormValue("IDs")
		ekx := strings.Split(IDsAsString, ",")
		if len(ekx) == 0 {
			log.Printf("Path: %s, Error: no key\n", s.R.URL.Path)
			http.Error(s.W, "no key", http.StatusBadRequest)
			return
		}
		err := page.DeleteMulti(s, ekx)
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
		ps, crsr, err := page.GetMulti(s, crsr, nil, nil)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.PrevPageURL = "/pages?c=" + crsr.String()
		rb.Result = ps
		// LastModifed AND Created ALSO SENDING WITH THEIR O VALUES, FIND A WAY TO
		// REMOVE THEM.
		api.WriteResponseJSON(s, rb)
	}
}
