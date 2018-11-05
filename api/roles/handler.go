/*
Package roles handles role requests.
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

// Handler handles "POST" and "DELETE" requests.
func Handler(s *session.Session) {
	u, uKey, err := user.GetViaEmail(s)
	if err != nil {
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
