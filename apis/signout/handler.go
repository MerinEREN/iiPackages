/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package signout

import (
	api "github.com/MerinEREN/iiPackages/apis"
	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
	"log"
	"net/http"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	URL, err := user.LogoutURL(ctx, "/")
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rb := new(api.ResponseBody)
	rb.Result = URL
	api.WriteResponse(w, r, rb)
}
