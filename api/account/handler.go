/*
Package account handles account requests.
*/
package account

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
)

// Handler returns and modifies account entities.
func Handler(s *session.Session) {
	k := new(datastore.Key)
	var err error
	ID := s.R.URL.Path[len("/accounts/"):]
	if ID == "" && s.R.Method != "GET" {
		log.Printf("Path: %s, Error: no user ID\n", s.R.URL.Path)
		http.Error(s.W, "No user ID", http.StatusBadRequest)
		return
	}
	if ID != "" {
		k, err = datastore.DecodeKey(ID)
		if err != nil {
			log.Printf("Page:%s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	switch s.R.Method {
	case "PUT":
		// Handle PUT requests
	case "DELETE":
		err = datastore.Delete(s.Ctx, k)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles GET requests
		acc := new(account.Account)
		if ID == "" {
			item, err := memcache.Get(s.Ctx, "acc")
			if err == nil {
				err = json.Unmarshal(item.Value, acc)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n",
						s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				ku := new(datastore.Key)
				item, err = memcache.Get(s.Ctx, "uKey")
				if err == nil {
					err = json.Unmarshal(item.Value, ku)
					if err != nil {
						log.Printf("Page:%s, Error: %v\n",
							s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					acc, err = account.Get(s.Ctx, ku.Parent())
					if err != nil {
						log.Printf("Page:%s, Error: %v\n",
							s.R.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
				} else {
					ku, err = user.GetKeyViaEmail(s)
					if err == datastore.Done {
						// IMPOSIBLE BUT !!!!!!!!!!!!!!!!!!!!!!!!!!
						log.Printf("Page:%s, Error: %v\n",
							s.R.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(s.W, err.Error(),
							http.StatusNoContent)
						return
					} else if err != nil {
						log.Printf("Page:%s, Error: %v\n",
							s.R.URL.Path, err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					} else {
						acc, err = account.Get(s.Ctx, ku.Parent())
						if err != nil {
							log.Printf("Page:%s, Error: %v\n",
								s.R.URL.Path, err)
							// ALSO LOG THIS WITH DATASTORE LOG
							http.Error(s.W, err.Error(),
								http.StatusInternalServerError)
							return
						}
					}
					bs, err := json.Marshal(ku)
					if err != nil {
						log.Printf("Page:%s, Error: %v\n",
							s.R.URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "uKey",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Page:%s, Error: %v\n",
							s.R.URL.Path, err)
					}
				}
				bs, err := json.Marshal(acc)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n",
						s.R.URL.Path, err)
				}
				item = &memcache.Item{
					Key:   "acc",
					Value: bs,
				}
				err = memcache.Add(s.Ctx, item)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n",
						s.R.URL.Path, err)
				}
			}
		} else {
			acc, err = account.Get(s.Ctx, k)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
		}
		s.W.Header().Set("Content-Type", "application/json")
		api.WriteResponseJSON(s, acc)
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
	// http.NotFound(w, r)
	// http.Redirect(w, r, "/MerinEREN", http.StatusFound)
}
