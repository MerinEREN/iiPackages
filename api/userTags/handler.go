/*
Package userTags "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package userTags

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/datastore/userTag"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Handler "Exported functions should have a comment"
// REMOVE USER TAG REQUEST HANDLERS TO THE userTag HANDLER AFTER REGEX ROUTING DONE !!!!!!!
func Handler(s *session.Session) {
	uID := strings.Split(s.R.URL.Path, "/")[2]
	if uID == "" {
		log.Printf("Path: %s, Error: no tag ID to delete\n", s.R.URL.Path)
		http.Error(s.W, "No tag ID to delete", http.StatusBadRequest)
		return
	}
	uKey, err := datastore.DecodeKey(uID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// Or the status code may be "StatusBadRequest"
		http.Error(s.W, err.Error(), http.StatusBadRequest)
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
		var tIDx []string
		err = json.Unmarshal(bs, &tIDx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		var kx []*datastore.Key
		var utx userTag.UserTags
		for _, v := range tIDx {
			k := datastore.NewKey(s.Ctx, "UserTag", v, 0, uKey)
			kx = append(kx, k)
			tKey, err := datastore.DecodeKey(v)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			ut := &userTag.UserTag{
				TagKey: tKey,
			}
			utx = append(utx, ut)
		}
		err = userTag.PutMulti(s.Ctx, kx, utx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	case "DELETE":
		tID := strings.Split(s.R.URL.Path, "/")[3]
		if tID == "" {
			log.Printf("Path: %s, Error: no tag ID to delete\n", s.R.URL.Path)
			http.Error(s.W, "No tag ID to delete", http.StatusBadRequest)
			return
		}
		k := datastore.NewKey(s.Ctx, "UserTag", tID, 0, uKey)
		err = userTag.Delete(s.Ctx, k)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		ts := make(tag.Tags)
		kx, err := userTag.GetKeysUserOrTag(s.Ctx, uKey)
		if err == datastore.Done {
			if kx != nil {
				ts, err = tag.GetMulti(s, kx)
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
		rb.Result = ts
		api.WriteResponseJSON(s, rb)
	}
}
