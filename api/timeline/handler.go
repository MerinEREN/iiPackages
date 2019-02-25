// Package timeline returns the logged user's timeline entities count.
package timeline

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/demand"
	// "github.com/MerinEREN/iiPackages/datastore/offer"
	"github.com/MerinEREN/iiPackages/datastore/servicePack"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/datastore/userTag"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
)

// Handler returns offers, demands, servicePacks and other necessary item counts
// via user or account (if user is admin) tag keys.
// First checks memcache for tag keys, if not present gets tag keys via request's uID which
// is encoded key.
func Handler(s *session.Session) {
	err := s.R.ParseForm()
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	var userTagKeys []*datastore.Key
	var accTagKeys []*datastore.Key
	var tagKeysQuery []*datastore.Key
	uKey, err := datastore.DecodeKey(s.R.Form.Get("uID"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	u := new(user.User)
	item, err := memcache.Get(s.Ctx, "u")
	if err == nil {
		err = json.Unmarshal(item.Value, u)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n",
				s.R.URL.Path, err)
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else {
		err = datastore.Get(s.Ctx, uKey, u)
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNoContent)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n",
				s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		} else {
			bs, err := json.Marshal(u)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
			item = &memcache.Item{
				Key:   "u",
				Value: bs,
			}
			err = memcache.Add(s.Ctx, item)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
		}
	}
	isAdmin, err := u.IsAdmin(s.Ctx)
	if isAdmin {
		item, err := memcache.Get(s.Ctx, "accTagKeys")
		if err == nil {
			err = json.Unmarshal(item.Value, &accTagKeys)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			tagKeysQuery = accTagKeys
		} else {
			uKeys, err := user.GetKeysByParentOrdered(s.Ctx, uKey.Parent())
			if err == datastore.Done {
				if len(uKeys) == 0 {
					log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", s.R.URL.Path, err)
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
			} else if err != nil {
				log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			for _, v := range uKeys {
				userTagKeys, err = userTag.GetKeysUserOrTag(s.Ctx, v)
				if err == datastore.Done {
					if len(userTagKeys) == 0 {
						log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
						s.W.WriteHeader(http.StatusNoContent)
						return
					}
				} else if err != nil {
					log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for _, v2 := range userTagKeys {
					absent := true
					for _, v3 := range accTagKeys {
						if *v3 == *v2 {
							absent = false
						}
					}
					if absent {
						accTagKeys = append(accTagKeys, v2)
					}
				}
			}
			bs, err := json.Marshal(accTagKeys)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
			item = &memcache.Item{
				Key:   "accTagKeys",
				Value: bs,
			}
			err = memcache.Add(s.Ctx, item)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
			tagKeysQuery = accTagKeys
		}
	} else {
		item, err := memcache.Get(s.Ctx, "userTagKeys")
		if err == nil {
			err = json.Unmarshal(item.Value, &tagKeysQuery)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			userTagKeys, err = userTag.GetKeysUserOrTag(s.Ctx, uKey)
			if err == datastore.Done {
				if len(userTagKeys) == 0 {
					log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
			} else if err != nil {
				log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			bs, err := json.Marshal(userTagKeys)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
			item = &memcache.Item{
				Key:   "userTagKeys",
				Value: bs,
			}
			err = memcache.Add(s.Ctx, item)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
			tagKeysQuery = userTagKeys
		}
	}
	var crsrAsStringDx []string
	if len(s.R.Form["cds"]) == 0 {
		crsrAsStringDx = make([]string, len(tagKeysQuery))
	} else {
		crsrAsStringDx = s.R.Form["cds"]
	}
	/* var crsrAsStringOx []string
	if len(s.R.Form["cos"]) == 0 {
		crsrAsStringOx = make([]string, len(tagKeysQuery))
	} else {
		crsrAsStringOx = s.R.Form["cos"]
	} */
	var crsrAsStringSPx []string
	if len(s.R.Form["csps"]) == 0 {
		crsrAsStringSPx = make([]string, len(tagKeysQuery))
	} else {
		crsrAsStringSPx = s.R.Form["csps"]
	}
	// var kdx, kox, kspx []*datastore.Key
	var kdx, kspx []*datastore.Key
	kds := make(map[string]*datastore.Key)
	ksps := make(map[string]*datastore.Key)
	// FIND A WAY TO MAKE QUERIES CONQURENTLY
	for i, v := range tagKeysQuery {
		kdx, err = demand.GetNewestKeysFilteredByTag(s.Ctx, crsrAsStringDx[i], v)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Request: get newest demands via tag key to count, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
		for _, v := range kdx {
			kds[v.Encode()] = v
		}
		/* kox, err = offer.GetNewestKeysFilteredByTag(s.Ctx, crsrAsStringOx[i], v)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Request: get newest offers via tag key to count, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
		for _, v := range kox {
			kos[v.Encode()] = v
		} */
		kspx, err = servicePack.GetNewestKeysFilteredByTag(s.Ctx, crsrAsStringSPx[i], v)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Request: get newest service packs via tag key to count, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
		for _, v := range kspx {
			ksps[v.Encode()] = v
		}
	}
	rb := new(api.ResponseBody)
	// rb.Result = len(kds) + len(kos) + len(ksps)
	rb.Result = len(kds) + len(ksps)
	if rb.Result == 0 {
		s.W.WriteHeader(http.StatusNoContent)
		return
	}
	api.WriteResponseJSON(s, rb)
}
