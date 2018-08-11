/*
Package offer "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package offer

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/offer"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/datastore/userTag"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
)

// Handler returns account offers via account ID if provided.
// Otherwise if the user is "admin" returns offers via account tag keys, else
// returns only via user tag keys.
// SAVE WITH ACCOUNT KEY AS THE PARENT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Handler(s *session.Session) {
	accID := s.R.FormValue("aID")
	crsr, err := datastore.DecodeCursor(s.R.FormValue("c"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	rb := new(api.ResponseBody)
	var os offer.Offers
	if accID != "" {
		aKey, err := datastore.DecodeKey(accID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		os, crsr, err = offer.GetViaParent(s.Ctx, crsr, aKey)
		if err == datastore.Done {
			if len(os) == 0 {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
		} else if err != nil {
			log.Printf("Path: %s, Request: get offers via parent key, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.PrevPageURL = "/offers?aID=" + accID + "&c=" + crsr.String()
	} else {
		// For timeline
		var userTagKeys []*datastore.Key
		var accTagKeys []*datastore.Key
		var tagKeysQuery []*datastore.Key
		drct := s.R.FormValue("d")
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
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
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
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			} else {
				bs, err := json.Marshal(u)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				}
				item = &memcache.Item{
					Key:   "u",
					Value: bs,
				}
				err = memcache.Add(s.Ctx, item)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				}
			}
		}
		if u.IsAdmin() {
			item, err := memcache.Get(s.Ctx, "accTagKeys")
			if err == nil {
				err = json.Unmarshal(item.Value, &accTagKeys)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
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
			item, err = memcache.Get(s.Ctx, "userTagKeys")
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
		switch drct {
		case "next":
			os, crsr, err = offer.GetNewest(s.Ctx, crsr, tagKeysQuery)
			if err == datastore.Done {
				if len(os) == 0 {
					log.Printf("Path: %s, Request: get newer offers, Error: %v\n", s.R.URL.Path, err)
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
			} else if err != nil {
				log.Printf("Path: %s, Request: get newer offers, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			rb.NextPageURL = "/offers?uID=" + s.R.FormValue("uID") + "&d=next&c=" + crsr.String()
		case "prev":
			os, crsr, err = offer.GetOldest(s.Ctx, crsr, tagKeysQuery)
			if err == datastore.Done {
				if len(os) == 0 {
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
			} else if err != nil {
				log.Printf("Path: %s, Request: get older offers, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			rb.PrevPageURL = "/offers?uID=" + s.R.FormValue("uID") + "&d=prev&c=" + crsr.String()
		default:
			var crsrNew datastore.Cursor
			var crsrOld datastore.Cursor
			os, crsrNew, crsrOld, err = offer.Get(s.Ctx, tagKeysQuery)
			if err == datastore.Done {
				if len(os) == 0 {
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
			} else if err != nil {
				log.Printf("Path: %s, Request: get offers, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			rb.NextPageURL = "/offers?uID=" + s.R.FormValue("uID") + "&d=next&c=" + crsrNew.String()
			rb.PrevPageURL = "/offers?uID=" + s.R.FormValue("uID") + "&d=prev&c=" + crsrOld.String()
		}
	}
	rb.Result = os
	api.WriteResponse(s, rb)
}
