/*
Package userRoles "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package userRoles

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/role"
	"github.com/MerinEREN/iiPackages/datastore/userRole"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Handler "Exported functions should have a comment"
// REMOVE USER ROLE DELETE REQUEST HANDLER TO THE userRole HANDLER AFTER REGEX ROUTING DONE
func Handler(s *session.Session) {
	uID := strings.Split(s.R.URL.Path, "/")[2]
	if uID == "" {
		log.Printf("Path: %s, Error: no user ID to delete\n", s.R.URL.Path)
		http.Error(s.W, "No user ID to delete", http.StatusBadRequest)
		return
	}
	uKey, err := datastore.DecodeKey(uID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// Or the status code may be "StatusBadRequest"
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	switch s.R.Method {
	case "POST":
		var bs []byte
		bs, err = ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		var rIDx []string
		err = json.Unmarshal(bs, &rIDx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		var kx []*datastore.Key
		var urx userRole.UserRoles
		for _, v := range rIDx {
			k := datastore.NewKey(s.Ctx, "UserRole", v, 0, uKey)
			kx = append(kx, k)
			rKey, err := datastore.DecodeKey(v)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			ur := &userRole.UserRole{
				RoleKey: rKey,
			}
			urx = append(urx, ur)
		}
		err = userRole.PutMulti(s.Ctx, kx, urx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	case "DELETE":
		rID := strings.Split(s.R.URL.Path, "/")[3]
		if rID == "" {
			log.Printf("Path: %s, Error: no role ID to delete\n", s.R.URL.Path)
			http.Error(s.W, "No role ID to delete", http.StatusBadRequest)
			return
		}
		k := datastore.NewKey(s.Ctx, "UserRole", rID, 0, uKey)
		err = userRole.Delete(s.Ctx, k)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		rs := make(role.Roles)
		kx, err := userRole.GetKeysUserOrRole(s.Ctx, uKey)
		if err == datastore.Done {
			if kx != nil {
				rs, err = role.GetMulti(s.Ctx, kx)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = rs
		api.WriteResponseJSON(s, rb)
	}
}
