/*
Package tagsDemand "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package tagsDemand

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/datastore/tagDemand"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Handler "Exported functions should have a comment"
// REMOVE DELETE REQUEST HANDLER TO THE "tagDemand" HANDLER AFTER REGEX ROUTING DONE !!!!!!
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	dID := q.Get("dID")
	if dID == "" {
		log.Printf("Path: %s, Error: No parent ID\n", URL.Path)
		http.Error(s.W, "No parent ID", http.StatusBadRequest)
		return
	}
	kd, err := datastore.DecodeKey(dID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	switch s.R.Method {
	case "POST":
		ct := s.R.Header.Get("Content-Type")
		if ct != "application/json" {
			log.Printf("Path: %s, Error: %v\n", URL.Path,
				"Content type is not application/json")
			http.Error(s.W, "Content type is not application/json",
				http.StatusUnsupportedMediaType)
			return
		}
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
		var tdx tagDemand.TagDemands
		for _, v := range tIDx {
			k := datastore.NewKey(s.Ctx, "TagDemand", v, 0, kd)
			kx = append(kx, k)
			kt, err := datastore.DecodeKey(v)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			td := &tagDemand.TagDemand{
				TagKey:  kt,
				Created: time.Now(),
			}
			tdx = append(tdx, td)
		}
		_, err = datastore.PutMulti(s.Ctx, kx, tdx)
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
		k := datastore.NewKey(s.Ctx, "TagDemand", tID, 0, kd)
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
		ktx, err := tagDemand.GetKeysByDemandOrTagKey(s.Ctx, kd)
		if err != nil {
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
