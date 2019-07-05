// Package signin returns the login url/urls.
package signin

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/session"
	googleUser "google.golang.org/appengine/user"
	"log"
	"net/http"
)

// Handler returns login urls.
func Handler(s *session.Session) {
	/* if s.R.URL.Path == "/favicon.ico" {
		return
	} */
	// For /favicon.ico
	if s.R.URL.Path != "/" {
		return
	}
	gURL, err := googleUser.LoginURL(s.Ctx, s.R.URL.String())
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	loginURLs := make(map[string]string)
	loginURLs["Google"] = gURL
	loginURLs["LinkedIn"] = gURL
	loginURLs["Twitter"] = gURL
	loginURLs["Facebook"] = gURL
	/* loginURLs["LinkedIn"] = "liURL"
	loginURLs["Twitter"] = "tURL"
	loginURLs["Facebook"] = "fURL" */
	s.W.Header().Set("Content-Type", "application/json")
	api.WriteResponseJSON(s, loginURLs)
}
