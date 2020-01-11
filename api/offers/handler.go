/*
Package offers "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package offers

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/offer"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handler returns account's offers via user ID if the user is admin
// otherwise returns only logged user's offers.
// Or returns demand's offers via demand ID.
// If the request method is POST, puts the offer to the datastore.
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	dID := q.Get("dID")
	switch s.R.Method {
	case "POST":
		if dID == "" {
			log.Printf("Path: %s, Error: no demand ID\n", URL.Path)
			http.Error(s.W, "No demand ID", http.StatusBadRequest)
			return
		}
		ct := s.R.Header.Get("Content-Type")
		if ct != "application/json" {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
				"Content type is not application/json")
			http.Error(s.W, "Content type is not application/json",
				http.StatusUnsupportedMediaType)
			return
		}
		bs, err := ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		o := new(offer.Offer)
		err = json.Unmarshal(bs, o)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		o.Status = "pending"
		o.Created = time.Now()
		pk, err := datastore.DecodeKey(dID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		k := datastore.NewIncompleteKey(s.Ctx, "Offer", pk)
		k, err = datastore.Put(s.Ctx, k, o)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		o.ID = k.Encode()
		sx := []string{"/offers", o.ID}
		path := strings.Join(sx, "/")
		rel, err := URL.Parse(path)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.Header().Set("Location", rel.String())
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponseJSON(s, o)
	default:
		accID := q.Get("aID")
		os := make(offer.Offers)
		if accID != "" {
			// UPDATE THIS BLOCK !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			ka, err := datastore.DecodeKey(accID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			kux, err := user.GetKeysByParentOrdered(s.Ctx, ka)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			after := q["after"]
			if len(after) == 0 {
				after = make([]string, len(kux))
			}
			var lim int
			limit := q.Get("limit")
			if limit == "" {
				lim = 0
			} else {
				lim, err = strconv.Atoi(limit)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				}
			}
			for i, v := range kux {
				ds2, crsrAsString, err := demand.GetNextByParentLimited(
					s.Ctx, after[i], v, lim)
				if err != nil {
					log.Printf("Path: %s, Request: get account demands via users keys, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i2, v2 := range ds2 {
					ds[i2] = v2
				}
				crsrAsStringx = append(crsrAsStringx, crsrAsString)
			}
			next := api.GenerateSubLink(s, crsrAsStringx, "next")
			s.W.Header().Set("Link", next)
		} else if dID != "" {
			kd, err := datastore.DecodeKey(dID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			crsrAsString := q.Get("after")
			var lim int
			limit := q.Get("limit")
			if limit == "" {
				lim = 0
			} else {
				lim, err = strconv.Atoi(limit)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				}
			}
			os, crsrAsString, err = offer.GetNextByParentLimited(s.Ctx, crsrAsString, kd, lim)
			if err != nil {
				log.Printf("Path: %s, Request: get demand offers via demand key, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			switch len(os) {
			case 0:
				s.W.WriteHeader(http.StatusNoContent)
			case 1:
				next := api.GenerateSubLink(s, crsrAsString, "next")
				s.W.Header().Set("Link", next)
				s.W.Header().Set("Content-Type", "application/json")
				for _, v := range os {
					api.WriteResponseJSON(s, v)
				}
			default:
				next := api.GenerateSubLink(s, crsrAsString, "next")
				s.W.Header().Set("Link", next)
				s.W.Header().Set("Content-Type", "application/json")
				api.WriteResponseJSON(s, os)
			}
		} else {
			// For timeline
		}
	}
}
