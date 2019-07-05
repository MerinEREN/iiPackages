/*
Package tagsUser "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package tagsUser

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/datastore/tagUser"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
)

// Handler "Exported functions should have a comment"
// REMOVE USER TAG DELETE REQUEST HANDLER TO THE tagUser HANDLER AFTER REGEX ROUTING DONE !
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	pID := q.Get("pID")
	if pID == "" {
		log.Printf("Path: %s, Error: No parent ID\n", URL.Path)
		http.Error(s.W, "No parent ID", http.StatusBadRequest)
		return
	}
	pk, err := datastore.DecodeKey(pID)
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
		var tIDx []string
		err = json.Unmarshal(bs, &tIDx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		var kx []*datastore.Key
		var tux tagUser.TagsUser
		for _, v := range tIDx {
			k := datastore.NewKey(s.Ctx, "TagUser", v, 0, pk)
			kx = append(kx, k)
			kt, err := datastore.DecodeKey(v)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			tu := &tagUser.TagUser{
				TagKey: kt,
			}
			tux = append(tux, tu)
		}
		_, err = datastore.PutMulti(s.Ctx, kx, tux)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	case "DELETE":
		tID := q.Get("tID")
		if tID == "" {
			log.Printf("Path: %s, Error: no tag ID to delete\n", URL.Path)
			http.Error(s.W, "No tag ID to delete", http.StatusBadRequest)
			return
		}
		k := datastore.NewKey(s.Ctx, "TagUser", tID, 0, pk)
		err = datastore.Delete(s.Ctx, k)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "GET" requests
		ts := make(tag.Tags)
		ktx, err := tagUser.GetKeysByUserOrTagKey(s.Ctx, pk)
		if err == datastore.Done {
			if ktx != nil {
				ts, err = tag.GetMulti(s.Ctx, ktx)
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
		s.W.Header().Set("Content-Type", "application/json")
		switch len(ts) {
		case 1:
			for _, v := range ts {
				api.WriteResponseJSON(s, v)
			}
		default:
			api.WriteResponseJSON(s, ts)
		}
	}
}
