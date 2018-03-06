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
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/servicePack"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
)

func Handler(s *session.Session) {
	var uTagIDs []*datastore.Key
	item, err := memcache.Get(s.Ctx, "uTagIDs")
	if err == nil {
		err = json.Unmarshal(item.Value, &uTagIDs)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n",
				s.R.URL.Path, err)
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
	} else {
		item, err = memcache.Get(s.Ctx, "u")
		if err == nil {
			u := new(user.User)
			err = json.Unmarshal(item.Value, u)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n",
					s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			uTagIDs = u.TagIDs
		} else {
			uTagIDs, err = user.GetTagIDs(s.Ctx, s.U.Email)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		bs, err := json.Marshal(uTagIDs)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		}
		item = &memcache.Item{
			Key:   "uTagIDs",
			Value: bs,
		}
		err = memcache.Add(s.Ctx, item)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		}
	}
	c, err := datastore.DecodeCursor(s.R.FormValue("c"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	rb := new(api.ResponseBody)
	var servicePacks servicePack.ServicePacks
	switch s.R.FormValue("d") {
	case "next":
		var cursor datastore.Cursor
		servicePacks, cursor, err = servicePack.GetNewest(s.Ctx, c, uTagIDs)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.NextPageURL = "/servicePacks?d=next&" + "c=" + cursor.String()
	case "prev":
		var cursor datastore.Cursor
		servicePacks, cursor, err = servicePack.GetOldest(s.Ctx, c, uTagIDs)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.PrevPageURL = "/servicePacks?d=prev&" + "c=" + cursor.String()
	default:
		var cNew datastore.Cursor
		var cOld datastore.Cursor
		servicePacks, cNew, cOld, err = servicePack.Get(s.Ctx, uTagIDs)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb.NextPageURL = "/servicePacks?d=next&" + "c=" + cNew.String()
		rb.PrevPageURL = "/servicePacks?d=prev&" + "c=" + cOld.String()
	}
	rb.Result = servicePacks
	api.WriteResponse(s, rb)
}
