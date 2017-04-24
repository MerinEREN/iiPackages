/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package logout

import (
	// "encoding/json"
	// "fmt"
	// "github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/cookie"
	// "github.com/MerinEREN/iiPackages/page/content"
	// "github.com/MerinEREN/iiPackages/page/template"
	// usr "github.com/MerinEREN/iiPackages/datastore/user"
	// "google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
	// "google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	// "html/template"
	"log"
	// "mime/multipart"
	"net/http"
	// "regexp"
	"golang.org/x/net/context"
	// "time"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	url, err := user.LogoutURL(ctx, "/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Deleting session cookie
	if err = cookie.Delete(w, r, "session"); err != nil {
		log.Printf("Path: %s. Error: %v\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//  CHANGE NECESSARY DB FIELDS OF USER !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	http.Redirect(w, r, url, http.StatusFound)
}
