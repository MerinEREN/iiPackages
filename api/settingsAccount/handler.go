// Package settingsAccount returns account struct.
package settingsAccount

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler returns account struct via encoded account key.
func Handler(s *session.Session) {
	acc := new(account.Account)
	k, err := datastore.DecodeKey(s.R.FormValue("aID"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusBadRequest)
		return
	}
	err = datastore.Get(s.Ctx, k, acc)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	rb := new(api.ResponseBody)
	rb.Result = acc
	api.WriteResponseJSON(s, rb)
}
