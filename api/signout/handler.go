// Package signout handles signout request.
package signout

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/user"
	"log"
	"net/http"
)

// Handler returns user's logout url as a json encoded data.
func Handler(s *session.Session) {
	URL, err := user.LogoutURL(s.Ctx, "/")
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	rb := new(api.ResponseBody)
	rb.Result = URL
	api.WriteResponseJSON(s, rb)
}
