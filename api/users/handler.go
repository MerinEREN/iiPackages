/*
Package users returns users with limited properties to list in User Settings page.
*/
package users

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler returns non admin users via account id.
func Handler(s *session.Session) {
	accKey, err := datastore.DecodeKey(s.R.FormValue("accID"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	var us user.Users
	us, err = user.GetProjected(s.Ctx, accKey)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	rb := new(api.ResponseBody)
	rb.Result = us
	api.WriteResponse(s, rb)
}
