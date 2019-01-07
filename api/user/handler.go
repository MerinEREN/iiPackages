// Package user deletes, updates and returns an user.
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

// Handler deletes, updates and returns user via user ID which is an encoded key
// if provided, otherwise returns logged user.
func Handler(s *session.Session) {
	rb := new(api.ResponseBody)
	ID := strings.Split(s.R.URL.Path, "/")[2]
	if ID == "" && s.R.Method != "GET" {
		log.Printf("Path: %s, Error: no user ID\n", s.R.URL.Path)
		http.Error(s.W, "No user ID", http.StatusBadRequest)
		return
	}
	uLogged := new(user.User)
	u := new(user.User)
	k := new(datastore.Key)
	var err error
	switch s.R.Method {
	case "DELETE":
		err = user.Delete(s.Ctx, ID)
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
		s.W.WriteHeader(http.StatusNoContent)
	case "PUT":
	default:
		if ID != "" {
			k, err = datastore.DecodeKey(ID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
				http.Error(s.W, err.Error(), http.StatusBadRequest)
				return
			}
			// Do not return logged user.
			uLogged, err = Get(s)
			if uLogged.Email == k.StringID() {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			u, err = user.Get(s.Ctx, k)
			if err == datastore.ErrNoSuchEntity {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
				http.Error(s.W, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			// Do not return deleted user.
			if u.Status == "deleted" {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
			us := make(user.Users)
			us[u.ID] = u
			rb.Result = us
		} else {
			// Return logged user.
			item, err := memcache.Get(s.Ctx, "u")
			if err == nil {
				err = json.Unmarshal(item.Value, u)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
						err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				k, err = datastore.DecodeKey(u.ID)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
						err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				var bs []byte
				item, err = memcache.Get(s.Ctx, "uKey")
				if err == nil {
					err = json.Unmarshal(item.Value, k)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					u, err = user.Get(s.Ctx, k)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
				} else {
					u, k, err = user.GetViaEmail(s)
					if err == datastore.Done {
						acc := new(account.Account)
						acc, u, k, err = user.CreateWithAccount(s)
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
					bs, err = json.Marshal(k)
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
			if u.Status == "deleted" || u.Status == "suspended" {
				s.W.WriteHeader(http.StatusUnauthorized)
				return
			}
			// USELESS FOR NOW, store some session data if needed !!!!!
			// If cookie present does nothing.
			// So does not necessary to check.
			err = cookie.Set(s, "session", s.U.Email)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
			rb.Result = u
		}
		api.WriteResponseJSON(s, rb)
	}
}
