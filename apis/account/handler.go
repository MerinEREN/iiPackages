/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package account

import (
	"encoding/json"
	// "fmt"
	api "github.com/MerinEREN/iiPackages/apis"
	"github.com/MerinEREN/iiPackages/datastore/account"
	// "github.com/MerinEREN/iiPackages/page/content"
	usr "github.com/MerinEREN/iiPackages/datastore/user"
	// "google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"
	// "google.golang.org/appengine/memcache"
	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
	// "io/ioutil"
	// "html/template"
	"log"
	// "mime/multipart"
	"net/http"
	// "regexp"
	// "time"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	wb := new(api.ResponseBody)
	accName := r.URL.Path[len("/accounts/"):]
	log.Printf("Selected account is: %s\n", accName)
	switch r.Method {
	case "GET":
		// acc := new(account.Account)
		/* u := new(usr.User)
		aKey := new(datastore.Key)
		uKey := new(datastore.Key)
		item, err := memcache.Get(ctx, "u")
		if err == nil {
			err = json.Unmarshal(item.Value, u)
			if err != nil {
				log.Printf("Page:%s, Error: %v\n", s, err)
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
					log.Printf("Error while creating "+
						"account: %v\n", err)
					// ALSO LOG THIS WITH DATASTORE LOG
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				}
				bs, err := json.Marshal(acc)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n",
						s, err)
				}
				item = &memcache.Item{
					Key:   "acc",
					Value: bs,
				}
				err = memcache.Add(ctx, item)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n",
						s, err)
				}
			case usr.ErrFindUser:
				log.Printf("Error while login user: %v\n",
					err)
				// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
				http.Error(w, err.Error(),
					http.StatusInternalServerError)
				return
			}
			bs, err := json.Marshal(u)
			if err != nil {
				log.Printf("Page:%s, Error: %v\n", s, err)
			}
			bsUKey, err := json.Marshal(uKey)
			if err != nil {
				log.Printf("Page:%s, Error: %v\n", s, err)
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
				log.Printf("Page:%s, Error: %v\n", s, err)
			}
		}
		item, err = memcache.Get(ctx, "acc")
		if err == nil {
			err = json.Unmarshal(item.Value, acc)
			if err != nil {
				log.Printf("Page:%s, Error: %v\n", s, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			item, err = memcache.Get(ctx, "uKey")
			if err == nil {
				err = json.Unmarshal(item.Value, uKey)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n", s, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				aKey = uKey.Parent()
				err = datastore.Get(ctx, aKey, acc)
				if err != nil {
					log.Printf("Error while getting user's "+
						"account data: %v\n", err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				uKey, err = usr.GetKey(ctx, ug.Email)
				switch err {
				case usr.ErrFindUser:
					log.Printf("Error while login user: %v\n",
						err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
					http.Error(w, err.Error(),
						http.StatusInternalServerError)
					return
				default:
					aKey = uKey.Parent()
					err = datastore.Get(ctx, aKey, acc)
					if err != nil {
						log.Printf("Error while getting user's "+
							"account data: %v\n", err)
						// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!
						http.Error(w, err.Error(),
							http.StatusInternalServerError)
						return
					}
				}
				bs, err := json.Marshal(uKey)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n", s, err)
				}
				item = &memcache.Item{
					Key:   "uKey",
					Value: bs,
				}
				err = memcache.Add(ctx, item)
				if err != nil {
					log.Printf("Page:%s, Error: %v\n", s, err)
				}
			}
			bs, err := json.Marshal(acc)
			if err != nil {
				log.Printf("Page:%s, Error: %v\n",
					s, err)
			}
			item = &memcache.Item{
				Key:   "acc",
				Value: bs,
			}
			err = memcache.Add(ctx, item)
			if err != nil {
				log.Printf("Page:%s, Error: %v\n",
					s, err)
			}
		}
		err = cookie.Set(w, r, "session", u.UUID)
		if err != nil {
			log.Printf("Error while creating session "+
				"cookie: %v\n", err)
		} */
		acc, err := account.Get(ctx, accName)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		au := struct {
			Account *account.Account `json:"account"`
			Users   *usr.Users       `json:"user"`
		}{acc, nil}
		wb.Result = au
	case "POST":
		// Handle POST requests
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
	wbm, err := json.Marshal(*wb)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(wbm)
	// http.NotFound(w, r)
	// http.Redirect(w, r, "/MerinEREN", http.StatusFound)
}
