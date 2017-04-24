/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package servicePack

import (
	"encoding/json"
	api "github.com/MerinEREN/iiPackages/apis"
	"github.com/MerinEREN/iiPackages/datastore/servicePack"
	usr "github.com/MerinEREN/iiPackages/datastore/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	"log"
	"net/http"
)

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, ug *user.User) {
	var uTagIDs []*datastore.Key
	item, err := memcache.Get(ctx, "uTagIDs")
	if err == nil {
		err = json.Unmarshal(item.Value, &uTagIDs)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n",
				r.URL.Path, err)
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else {
		item, err = memcache.Get(ctx, "u")
		if err == nil {
			u := new(usr.User)
			err = json.Unmarshal(item.Value, u)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n",
					r.URL.Path, err)
				http.Error(w, err.Error(),
					http.StatusInternalServerError)
				return
			}
			uTagIDs = u.TagIDs
		} else {
			uTagIDs, err = usr.GetTagIDs(ctx, ug.Email)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		bs, err := json.Marshal(uTagIDs)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		}
		item = &memcache.Item{
			Key:   "uTagIDs",
			Value: bs,
		}
		err = memcache.Add(ctx, item)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		}
	}
	c, err := datastore.DecodeCursor(r.FormValue("c"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
	}
	rb := new(api.ResponseBody)
	var servicePacks []*servicePack.ServicePack
	var cursor datastore.Cursor
	switch r.FormValue("d") {
	case "next":
		servicePacks, cursor, err = servicePack.GetNewest(ctx, c, uTagIDs)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		}
	case "prev":
		servicePacks, cursor, err = servicePack.GetOldest(ctx, c, uTagIDs)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		}
	default:
		servicePacks, cursor, err = servicePack.Get(ctx, uTagIDs)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		}
	}
	rb.Result = servicePacks
	rb.NextPageUrl = "/servicePacks?d=" + r.FormValue("d") + "&" + "c=" +
		cursor.String()
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}
