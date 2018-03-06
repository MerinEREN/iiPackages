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
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
)

func Handler(s *session.Session) {
	u, _, err := user.GetWithEmail(s.Ctx, s.U.Email)
	if err == user.ErrFindUser {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status == "frozen" {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		// fmt.Fprintf(s.W, "Permission denied !!!")
		http.Error(s.W, "some error message", http.StatusForbidden)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	// Always send corresponding header values instead of defaults !!!!
	s.W.Write(b)
}
