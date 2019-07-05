// Package timeline returns the logged user's timeline entities count.
package timeline

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/tagDemand"
	// "github.com/MerinEREN/iiPackages/datastore/offer"
	"github.com/MerinEREN/iiPackages/datastore/tagServicePack"
	"github.com/MerinEREN/iiPackages/datastore/tagUser"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
)

// Handler returns demands, offers, servicePacks and other necessary timeline item counts
// via user or account (if user is admin) tag keys.
// First checks memcache for tag keys
// if not present gets tag keys via request's uID which is encoded key.
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	uID := q.Get("uID")
	if uID == "" {
		log.Printf("Path: %s, Error: No user ID.\n", URL.Path)
		http.Error(s.W, "No user ID.", http.StatusBadRequest)
		return
	}
	ku, err := datastore.DecodeKey(uID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	u := new(user.User)
	item, err := memcache.Get(s.Ctx, "u")
	if err == nil {
		err = json.Unmarshal(item.Value, u)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = datastore.Get(s.Ctx, ku, u)
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNoContent)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		} else {
			bs, err := json.Marshal(u)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			}
			item = &memcache.Item{
				Key:   "u",
				Value: bs,
			}
			err = memcache.Add(s.Ctx, item)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			}
		}
	}
	var ktux []*datastore.Key
	var ktax []*datastore.Key
	var ktx []*datastore.Key
	isAdmin, err := u.IsAdmin(s.Ctx)
	if isAdmin {
		item, err := memcache.Get(s.Ctx, "ktax")
		if err == nil {
			err = json.Unmarshal(item.Value, &ktax)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			ktx = ktax
		} else {
			kux, err := user.GetKeysByParentOrdered(s.Ctx, ku.Parent())
			if err != nil {
				log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(kux) == 0 {
				// Impossible !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", URL.Path, err)
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
			for _, v := range kux {
				ktux, err = tagUser.GetKeysByUserOrTagKey(s.Ctx, v)
				if err == datastore.Done {
					if len(ktux) == 0 {
						log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
						s.W.WriteHeader(http.StatusNoContent)
						return
					}
				} else if err != nil {
					log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for _, v2 := range ktux {
					absent := true
					for _, v3 := range ktax {
						if *v3 == *v2 {
							absent = false
							break
						}
					}
					if absent {
						ktax = append(ktax, v2)
					}
				}
			}
			bs, err := json.Marshal(ktax)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			}
			item = &memcache.Item{
				Key:   "ktax",
				Value: bs,
			}
			err = memcache.Add(s.Ctx, item)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			}
			ktx = ktax
		}
	} else {
		item, err := memcache.Get(s.Ctx, "ktux")
		if err == nil {
			err = json.Unmarshal(item.Value, &ktx)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			ktux, err = tagUser.GetKeysByUserOrTagKey(s.Ctx, ku)
			if err == datastore.Done {
				if len(ktux) == 0 {
					log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
			} else if err != nil {
				log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			bs, err := json.Marshal(ktux)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			}
			item = &memcache.Item{
				Key:   "ktux",
				Value: bs,
			}
			err = memcache.Add(s.Ctx, item)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			}
			ktx = ktux
		}
	}
	befored := q["befored"]
	// beforeo := q["beforeo"]
	beforesp := q["beforesp"]
	kds := make(map[string]*datastore.Key)
	ksps := make(map[string]*datastore.Key)
	// FIND A WAY TO MAKE QUERIES CONQURENTLY
	for i, v := range ktx {
		kdx, _, err := tagDemand.GetPrevKeysParentsFilteredByTagKey(s.Ctx,
			befored[i], v)
		if err != datastore.Done {
			log.Printf("Path: %s, Request: get previous demand's keys by tag key to count, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v2 := range kdx {
			kds[v2.Encode()] = v2
		}
		/* kox, _, err := tagOffer.GetPrevKeysParentsFilteredByTagKey(s.Ctx, beforeo[i], v)
		if err != datastore.Done {
			log.Printf("Path: %s, Request: get previous offer's keys by tag key to count, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v3 := range kox {
			kos[v3.Encode()] = v3
		} */
		kspx, _, err := tagServicePack.GetPrevKeysParentsFilteredByTagKey(s.Ctx,
			beforesp[i], v)
		if err != datastore.Done {
			log.Printf("Path: %s, Request: get previous service pack's keys by tag key to count, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v4 := range kspx {
			ksps[v4.Encode()] = v4
		}
	}
	// count := len(kds) + len(kos) + len(ksps)
	count := len(kds) + len(ksps)
	if count == 0 {
		s.W.WriteHeader(http.StatusNoContent)
		return
	}
	s.W.Header().Set("Content-Type", "application/json")
	api.WriteResponseJSON(s, count)
}
