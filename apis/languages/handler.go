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
	api "github.com/MerinEREN/iiPackages/apis"
	"github.com/MerinEREN/iiPackages/datastore/language"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler "Exported functions should have a comment"
func Handler(s *session.Session) {
	rb := new(api.ResponseBody)
	switch s.R.Method {
	case "POST":
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
		// Reset the cursor.
		var c datastore.Cursor
		var langs language.Languages
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) error {
			var err1 error
			langs, c, err1 = language.PutAndGetMulti(s, c, lang)
			return err1
		}, nil)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.Result = langs
		rb.Reset = true
		rb.PrevPageURL = "/languages?d=prev&" + "c=" + c.String()
		s.W.WriteHeader(http.StatusCreated)
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
			langs, c, err := language.GetMulti(s, c)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			rb.PrevPageURL = "/languages?d=prev&" + "c=" + c.String()
			rb.Result = langs
		}
	}
	api.WriteResponse(s, rb)
}
