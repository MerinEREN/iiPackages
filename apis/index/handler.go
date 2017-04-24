/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package index

import (
	"encoding/json"
	api "github.com/MerinEREN/iiPackages/apis"
	"github.com/MerinEREN/iiPackages/cookie"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	usr "github.com/MerinEREN/iiPackages/datastore/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	"log"
	"net/http"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	if r.URL.Path == "/favicon.ico" {
		return
	}
	var bs []byte
	switch r.Method {
	case "GET":
		resBody := new(api.ResponseBody)
		// Login or get data needed
		if ug == nil {
			gURL, err := user.LoginURL(ctx, r.URL.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			loginURLs := make(map[string]string)
			loginURLs["Google"] = gURL
			loginURLs["LinkedIn"] = gURL
			loginURLs["Twitter"] = gURL
			loginURLs["Facebook"] = gURL
			resBody.Result = loginURLs
			// Also send general statistics data.
		} else {
			acc := new(account.Account)
			u := new(usr.User)
			// Users own account page or not
			// if ac.ID, ok := reqBodyUm["data"]["acc"]; !ok {
			aKey := new(datastore.Key)
			uKey := new(datastore.Key)
			p := new(photo.Photo)
			item, err := memcache.Get(ctx, "u")
			if err == nil {
				err = json.Unmarshal(item.Value, u)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				u, uKey, err = usr.Get(ctx, ug.Email)
				switch err {
				case datastore.Done:
					acc, u, uKey, err = account.Create(ctx)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG
						http.Error(w, err.Error(),
							http.StatusInternalServerError)
						return
					}
					bs, err = json.Marshal(acc)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "acc",
						Value: bs,
					}
					err = memcache.Add(ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
					}
				case usr.ErrFindUser:
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				}
				p, _, err = photo.Get(ctx, uKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
				} else {
					u.Photo = *p
				}
				bs, err = json.Marshal(u)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
				}
				bsUKey, err := json.Marshal(uKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
				}
				items := []*memcache.Item{
					{
						Key:   "u",
						Value: bs,
					},
					{
						Key:   "uKey",
						Value: bsUKey,
					},
				}
				err = memcache.AddMulti(ctx, items)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
				}
			}
			item, err = memcache.Get(ctx, "acc")
			if err == nil {
				err = json.Unmarshal(item.Value, acc)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				item, err = memcache.Get(ctx, "uKey")
				if err == nil {
					err = json.Unmarshal(item.Value, uKey)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
						http.Error(w, err.Error(),
							http.StatusInternalServerError)
						return
					}
					aKey = uKey.Parent()
					acc, err = account.Get(ctx, aKey)
					if err != nil && err != datastore.Done {
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(w, err.Error(),
							http.StatusInternalServerError)
						return
					}
				} else {
					uKey, err = usr.GetKey(ctx, ug.Email)
					switch err {
					// Imposible case
					/* case datastore.Done:
					acc, u, uKey, err = account.Create(ctx)
					if err != nil {
						log.Printf("Error while creating "+
							"account: %v\n", err)
						// ALSO LOG THIS WITH DATASTORE LOG
						http.Error(w, err.Error(),
							http.StatusInternalServerError)
						return
					}
					bs, err := json.Marshal(acc)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s, err)
					}
					itemA = &memcache.Item{
						Key:   "acc",
						Value: bs,
					}
					err = memcache.Add(ctx, itemA)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s, err)
					} */
					case usr.ErrFindUser:
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(w, err.Error(),
							http.StatusInternalServerError)
						return
					default:
						aKey = uKey.Parent()
						acc, err = account.Get(ctx, aKey)
						if err != nil && err != datastore.Done {
							log.Printf("Path: %s, Error: %v\n",
								r.URL.Path, err)
							// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
							http.Error(w, err.Error(),
								http.StatusInternalServerError)
							return
						}
					}
					bs, err = json.Marshal(uKey)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "uKey",
						Value: bs,
					}
					err = memcache.Add(ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							r.URL.Path, err)
					}
				}
				p, _, err = photo.Get(ctx, aKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
				} else {
					acc.Photo = *p
				}
				bs, err = json.Marshal(acc)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
				}
				item = &memcache.Item{
					Key:   "acc",
					Value: bs,
				}
				err = memcache.Add(ctx, item)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						r.URL.Path, err)
				}
			}
			// If cookie present does nothing.
			// So does not necessary to check.
			err = cookie.Set(w, r, "session", u.ID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n",
					r.URL.Path, err)
			}
			/* } else {
				// someone elses account
				s, ok := ac.ID.(string)
				if ok {
					acc, err = account.Get(ctx, s)
				} else {
					log.Println("Account name type is not string.")
					http.Error(w, "Account name type is not string.",
						http.StatusBadRequest)
					return
				}
			} */
			var ua userAccount
			ua.User = make(map[string]*usr.User)
			ua.User[u.ID] = u
			ua.Account = make(map[string]*account.Account)
			ua.Account[acc.ID] = acc
			resBody.Result = ua
		}
		/* t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		c := &http.Client{Transport: t}
		res, err := c.Get("file:///etc/passwd")
		log.Println(res, err) */
		// To respond to request without any data
		// w.WriteHeader(StatusOK)
		// Always send corresponding header values instead of defaults !!!!
		//w.Header().Set("Content-Type", "application/json; charset=utf-8")
		bs, err := json.Marshal(*resBody)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bs)
	case "POST":
		// Handle POST requests.
		// Allways close the body
		// defer r.Body.Close()
		// r.Body is io.ReadCloser type, so may be closing request body
		// explicitly is not necessary.
		/* bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var reqBodyUm map[string]map[string]interface{}
		err = json.Unmarshal(bs, &reqBodyUm)
		if err != nil {
			log.Println("Error while unmarshalling request body:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} */
		// Can use belowe line instead of ioutil.ReadAll() and json.Unmarshall()
		// But performs a little bit slower.
		// err = json.NewDecoder(r.Body()).Decode(&reqBodyUm)
	default:
		// Some default shit
	}
}
