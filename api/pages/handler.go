/*
Package pages "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package pages

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/page"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler "Exported functions should have a comment"
func Handler(s *session.Session) {
	rb := new(api.ResponseBody)
	switch s.R.Method {
	case "POST":
		keyName := s.R.FormValue("ID")
		title := s.R.FormValue("title")
		mpf, hdr, err := s.R.FormFile("file")
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		defer mpf.Close()
		p := &page.Page{
			ID:    keyName,
			Title: title,
			Mpf:   mpf,
			Hdr:   hdr,
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
		p.Link, err = storage.UploadFile(s, p.Mpf, p.Hdr)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		// Reset the cursor and get the entities from the begining.
		var c datastore.Cursor
		var ps page.Pages
		ps, c, err = page.PutAndGetMulti(s, c, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.Result = ps
		rb.Reset = true
		rb.PrevPageURL = "/pages?d=prev&" + "c=" + c.String()
		s.W.WriteHeader(http.StatusCreated)
	default:
		// Handles "GET" requests
		c, err := datastore.DecodeCursor(s.R.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		ps, c, err := page.GetMulti(s, c, nil, nil)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.PrevPageURL = "/pages?c=" + c.String()
		rb.Result = ps
		// LastModifed AND Created ALSO SENDING WITH THEIR O VALUES, FIND A WAY TO
		// REMOVE THEM.
	}
	api.WriteResponse(s, rb)
}
