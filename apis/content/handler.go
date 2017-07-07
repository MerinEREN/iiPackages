/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
// GET LANGUAGE CODE AND PAGE NAME FROM URL AND HANDLE THE REQUEST HERE.
package content

import (
	api "github.com/MerinEREN/iiPackages/apis"
	"github.com/MerinEREN/iiPackages/datastore/content"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

func Handler(s *session.Session) {
	if s.R.Method == "POST" {
		bs, err := ioutil.ReadAll(s.R.Body)
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
		}
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
		p, err = page.Put(s.Ctx, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(201)
		return
	} else {
		c, err := datastore.DecodeCursor(s.R.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		pages, c, err := page.GetMulti(s.Ctx, c)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.PrevPageURL = "/pages?d=prev&" + "c=" + c.String()
		rb.Result = pages
		api.WriteResponse(s, rb)
	}
}
