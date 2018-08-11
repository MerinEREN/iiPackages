/*
Package user returns non admin user when admin user selects from user settings.
*/
package user

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/cookie"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
	"strings"
)

// Handler returns user via user ID which is an encoded key if provided
// Otherwise returns logged user.
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	u := new(user.User)
	uKey := new(datastore.Key)
	var err error
	if ID != "" {
		uKey, err = datastore.DecodeKey(ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		err = datastore.Get(s.Ctx, uKey, u)
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Return logged user.
		item, err := memcache.Get(s.Ctx, "u")
		if err == nil {
			err = json.Unmarshal(item.Value, u)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			uKey, err = datastore.DecodeKey(u.ID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
		} else {
			var bs []byte
			item, err = memcache.Get(s.Ctx, "uKey")
			if err == nil {
				err = json.Unmarshal(item.Value, uKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				u, err = user.Get(s, uKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				u, uKey, err = user.GetViaEmail(s)
				if err == datastore.Done {
					acc := new(account.Account)
					acc, u, uKey, err = user.CreateWithAccount(s)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					bs, err = json.Marshal(acc)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "acc",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
				} else if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				bs, err = json.Marshal(uKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				}
				item = &memcache.Item{
					Key:   "uKey",
					Value: bs,
				}
				err = memcache.Add(s.Ctx, item)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				}
			}
			bs, err = json.Marshal(u)
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
		// USELESS FOR NOW, store some session data if needed !!!!!
		// If cookie present does nothing.
		// So does not necessary to check.
		err = cookie.Set(s, "session", s.U.Email)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		}
	}
	rb := new(api.ResponseBody)
	/* uKeys, err := user.GetUsersKeysViaParent(s.Ctx, uKey.Parent())
	if err != datastore.Done {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	tagKeys, err := userTags.Get(s.Ctx, uKeys)
	if err != nil && err != datastore.Done {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
		http.Error(s.W, err.Error(),
			http.StatusInternalServerError)
		return
	}
	tags, err := tags.GetMulti(s.Ctx, tagKeys)
	if err != nil && err != datastore.Done {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
		http.Error(s.W, err.Error(),
			http.StatusInternalServerError)
		return
	}
	u.Tags = tags*/
	rb.Result = u
	api.WriteResponse(s, rb)
}
