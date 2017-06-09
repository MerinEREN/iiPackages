/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package page

import (
	"encoding/json"
	api "github.com/MerinEREN/iiPackages/apis"
	"github.com/MerinEREN/iiPackages/datastore/page"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"io/ioutil"
	"log"
	"net/http"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	if r.Method == "POST" {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p := new(page.Page)
		err = json.Unmarshal(bs, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Using 'decoder' is an alternative and can be used if response body has
		// more than one json object.
		// Otherwise don't use it, because it has performance disadvantages
		// compared to first solution.
		/*decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} */
		p, err = page.Put(ctx, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(201)
		return
	} else {
		c, err := datastore.DecodeCursor(r.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pages, c, err := page.GetMulti(ctx, c)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.PrevPageURL = "/pages?d=prev&" + "c=" + c.String()
		rb.Result = pages
		api.WriteResponse(w, r, rb)
	}
}
