/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package roles

import (
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

func Handler(s *session.Session) {
	u, uKey, err := user.GetWithEmail(s.Ctx, s.U.Email)
	if err == user.ErrFindUser {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status != "suspended" {
		log.Printf("Suspended user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are suspended", http.StatusForbidden)
		return
	}
	if u.Type == "inHouse" && (u.IsAdmin() || u.IsContentEditor()) {
		acc := new(account.Account)
		aKey := uKey.Parent()
		err = datastore.Get(s.Ctx, aKey, acc)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are unauthorized user.", http.StatusUnauthorized)
		return
	}
}
