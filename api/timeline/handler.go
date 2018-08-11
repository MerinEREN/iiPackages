/*
Package timeline "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package timeline

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/datastore/offer"
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
	encCrsrD := s.R.FormValue("cd")
	encCrsrO := s.R.FormValue("co")
	encCrsrSp := s.R.FormValue("csp")
	var userTagKeys []*datastore.Key
	var accTagKeys []*datastore.Key
	var tagKeysQuery []*datastore.Key
	uKey, err := datastore.DecodeKey(s.R.FormValue("uID"))
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
	if u.IsAdmin() {
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
			uKeys, err := user.GetUsersKeysViaParent(s.Ctx, uKey.Parent())
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
				userTagKeys, err = userTag.Get(s.Ctx, v)
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
					accTagKeys = append(accTagKeys, v2)
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
			userTagKeys, err = userTag.Get(s.Ctx, uKey)
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
	rb := new(api.ResponseBody)
	// FIND A WAY TO MAKE QUERIES CONQURENTLY
	crsrD, err := datastore.DecodeCursor(encCrsrD)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	crsrO, err := datastore.DecodeCursor(encCrsrO)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	crsrSp, err := datastore.DecodeCursor(encCrsrSp)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	var countT, countD, countO, countSp int
	countD, err = demand.GetNewestCount(s.Ctx, crsrD, tagKeysQuery)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	countO, err = offer.GetNewestCount(s.Ctx, crsrO, tagKeysQuery)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	countSp, err = servicePack.GetNewestCount(s.Ctx, crsrSp, tagKeysQuery)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	countT = countD + countO + countSp
	rb.Result = countT
	api.WriteResponse(s, rb)
}
