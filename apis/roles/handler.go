/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package roles

import (
	// "encoding/json"
	// "fmt"
	"github.com/MerinEREN/iiPackages/datastore/account"
	usr "github.com/MerinEREN/iiPackages/datastore/user"
	// "google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	"golang.org/x/net/context"
	"log"
	// "mime/multipart"
	"net/http"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	u, uKey, err := usr.Get(ctx, ug.Email)
	if err == usr.ErrFindUser {
		log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status != "suspended" {
		log.Printf("Suspended user %s trying to see "+
			"%s path!!!", u.Email, r.URL.Path)
		http.Error(w, "You are suspended", http.StatusForbidden)
		return
	}
	if u.Type == "inHouse" && (u.IsAdmin() || u.IsContentEditor()) {
		acc := new(account.Account)
		aKey := uKey.Parent()
		err = datastore.Get(ctx, aKey, acc)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, r.URL.Path)
		http.Error(w, "You are unauthorized user.", http.StatusUnauthorized)
		return
	}
}
