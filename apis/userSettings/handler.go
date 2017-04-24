/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package userSettings

import (
	"encoding/json"
	// "fmt"
	// "github.com/MerinEREN/iiPackages/datastore/account"
	usr "github.com/MerinEREN/iiPackages/datastore/user"
	// "google.golang.org/appengine"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	"golang.org/x/net/context"
	"log"
	// "mime/multipart"
	"net/http"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	u, _, err := usr.Get(ctx, ug.Email)
	if err == usr.ErrFindUser {
		log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status == "frozen" {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, r.URL.Path)
		// fmt.Fprintf(w, "Permission denied !!!")
		http.Error(w, "some error message", http.StatusForbidden)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Always send corresponding header values instead of defaults !!!!
	w.Write(b)
}
