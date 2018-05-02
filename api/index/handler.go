/*
Package index "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package index

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/cookie"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	googleUser "google.golang.org/appengine/user"
	"log"
	"net/http"
)

// Handler "Exported functions should have a comment"
func Handler(s *session.Session) {
	/* if s.R.URL.Path == "/favicon.ico" {
		return
	} */
	if s.R.URL.Path != "/" {
		return
	}
	var bs []byte
	switch s.R.Method {
	case "POST":
		// Handle POST requests.
		// Allways close the body
		// defer r.Body.Close()
		// r.Body is io.ReadCloser type, so may be closing request body
		// explicitly is not necessary.
		/* bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(s.W, err.Error(), http.StatusBadRequest)
			return
		}
		var reqBodyUm map[string]map[string]interface{}
		err = json.Unmarshal(bs, &reqBodyUm)
		if err != nil {
			log.Println("Error while unmarshalling request body:", err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		} */
		// Can use belowe line instead of ioutil.ReadAll() and json.Unmarshall()
		// But performs a little bit slower.
		// err = json.NewDecoder(r.Body()).Decode(&reqBodyUm)
	default:
		// Handles "GET" requests
		rb := new(api.ResponseBody)
		// Login or get data needed
		if s.U == nil {
			googleURL, err := googleUser.LoginURL(s.Ctx, s.R.URL.String())
			if err != nil {
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			loginURLs := make(map[string]string)
			loginURLs["Google"] = googleURL
			loginURLs["LinkedIn"] = googleURL
			loginURLs["Twitter"] = googleURL
			loginURLs["Facebook"] = googleURL
			rb.Result = loginURLs
			// Also send general statistics data.
		} else {
			acc := new(account.Account)
			u := new(user.User)
			aKey := new(datastore.Key)
			uKey := new(datastore.Key)
			p := new(photo.Photo)
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
				u, uKey, err = user.GetWithEmail(s.Ctx, s.U.Email)
				switch err {
				case datastore.Done:
					acc, u, uKey, err = account.
						CreateAccountAndUser(s.Ctx)
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
				case user.ErrFindUser:
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				p, _, err = photo.Get(s.Ctx, uKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				} else {
					u.Photo = *p
				}
				bs, err = json.Marshal(u)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				}
				bsUKey, err := json.Marshal(uKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
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
				err = memcache.AddMulti(s.Ctx, items)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				}
			}
			item, err = memcache.Get(s.Ctx, "acc")
			if err == nil {
				err = json.Unmarshal(item.Value, acc)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
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
					aKey = uKey.Parent()
					acc, err = account.Get(s.Ctx, aKey)
					if err != nil && err != datastore.Done {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
				} else {
					uKey, err = user.GetKey(s.Ctx, s.U.Email)
					switch err {
					// Imposible case
					/* case datastore.Done:
					acc, u, uKey, err = account.Create(s.Ctx)
					if err != nil {
						log.Printf("Error while creating "+
							"account: %v\n", err)
						// ALSO LOG THIS WITH DATASTORE LOG
						http.Error(s.W, err.Error(),
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
					err = memcache.Add(s.Ctx, itemA)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s, err)
					} */
					case user.ErrFindUser:
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					default:
						aKey = uKey.Parent()
						acc, err = account.Get(s.Ctx, aKey)
						if err != nil && err != datastore.Done {
							log.Printf("Path: %s, Error: %v\n",
								s.R.URL.Path, err)
							// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
							http.Error(s.W, err.Error(),
								http.StatusInternalServerError)
							return
						}
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
				p, _, err = photo.Get(s.Ctx, aKey)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
				} else {
					acc.Photo = *p
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
			}
			// If cookie present does nothing.
			// So does not necessary to check.
			err = cookie.Set(s, "session", u.ID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n",
					s.R.URL.Path, err)
			}
			/* } else {
				// someone elses account
				s, ok := ac.ID.(string)
				if ok {
					acc, err = account.Get(s.Ctx, s)
				} else {
					log.Println("Account name type is not string.")
					http.Error(s.W, "Account name type is not string.",
						http.StatusBadRequest)
					return
				}
			} */
			var ua userAccount
			ua.User = u
			ua.Account = acc
			rb.Result = ua
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
		api.WriteResponse(s, rb)
	}
}
