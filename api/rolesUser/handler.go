/*
Package rolesUser "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package rolesUser

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/role"
	"github.com/MerinEREN/iiPackages/datastore/roleUser"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
)

// Handler "Exported functions should have a comment"
// REMOVE USER ROLE DELETE REQUEST HANDLER TO THE roleUser HANDLER AFTER REGEX ROUTING DONE
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	uID := q.Get("uID")
	if uID == "" {
		log.Printf("Path: %s, Error: No parent ID\n", URL.Path)
		http.Error(s.W, "No parent ID", http.StatusBadRequest)
		return
	}
	ku, err := datastore.DecodeKey(uID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	switch s.R.Method {
	case "POST":
		var bs []byte
		bs, err = ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		var rIDx []string
		err = json.Unmarshal(bs, &rIDx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		var kx []*datastore.Key
		var rux roleUser.RolesUser
		for _, v := range rIDx {
			k := datastore.NewKey(s.Ctx, "RoleUser", v, 0, ku)
			kx = append(kx, k)
			kr, err := datastore.DecodeKey(v)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			ru := &roleUser.RoleUser{
				RoleKey: kr,
			}
			rux = append(rux, ru)
		}
		_, err = datastore.PutMulti(s.Ctx, kx, rux)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	case "DELETE":
		rID := q.Get("rID")
		if rID == "" {
			log.Printf("Path: %s, Error: no role ID to delete\n", URL.Path)
			http.Error(s.W, "No role ID to delete", http.StatusBadRequest)
			return
		}
		k := datastore.NewKey(s.Ctx, "RoleUser", rID, 0, ku)
		err = datastore.Delete(s.Ctx, k)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		rs := make(role.Roles)
		krx, err := roleUser.GetKeysByUserOrRoleKey(s.Ctx, ku)
		if err == datastore.Done {
			if krx != nil {
				rs, err = role.GetMulti(s.Ctx, krx)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		switch len(rs) {
		case 1:
			s.W.Header().Set("Content-Type", "application/json")
			for _, v := range rs {
				api.WriteResponseJSON(s, v)
			}
		default:
			s.W.Header().Set("Content-Type", "application/json")
			api.WriteResponseJSON(s, rs)
		}
	}
}
