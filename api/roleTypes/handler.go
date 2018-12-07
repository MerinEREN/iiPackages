// Package roleTypes posts a roleType and also gets all roleTypes.
package roleTypes

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/roleType"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler posts a roleType and returns all the roleTypes from the begining of the kind
// and gets all the roleTypes from the begining of the kind.
func Handler(s *session.Session) {
	// THE CONTROLS BELOVE PREVENT GET REQUEST THAT NECESSARY FOR SELECT SOME SELECT
	// FIELDS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* u, err := user.Get(s)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status == "suspended" {
		log.Printf("Suspended user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are suspended", http.StatusForbidden)
		return
	}
	isAdmin, err := u.IsAdmin(s.Ctx)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	isContentEditor, err := u.IsContentEditor(s.Ctx)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Type != "inHouse" || !(isAdmin || isContentEditor) {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are unauthorized user.", http.StatusUnauthorized)
		return
	} */
	switch s.R.Method {
	case "POST":
		ID := s.R.FormValue("ID")
		rt := &roleType.RoleType{
			ID: ID,
		}
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		rts, err := roleType.PutAndGetMulti(s, rt)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = rts
		rb.PrevPageURL = "/roleTypes?c=" + crsr.String()
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	default:
		// Handles "GET" requests
		rts, err := roleType.GetMulti(s.Ctx, nil)
		if err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(rts) == 0 {
			s.W.WriteHeader(http.StatusNoContent)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = rts
		api.WriteResponseJSON(s, rb)
	}
}
