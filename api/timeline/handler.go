/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package timeline

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/datastore/offer"
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
	// FIND A WAY TO MAKE QUERIES CONQURENTLY
	cd, err := datastore.DecodeCursor(s.R.FormValue("cd"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	co, err := datastore.DecodeCursor(s.R.FormValue("co"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	csp, err := datastore.DecodeCursor(s.R.FormValue("csp"))
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	var countT, countD, countO, countSp int
	countD, err = demand.GetNewestCount(s.Ctx, cd, uTagIDs)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	countO, err = offer.GetNewestCount(s.Ctx, co, uTagIDs)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	countSp, err = servicePack.GetNewestCount(s.Ctx, csp, uTagIDs)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
	}
	countT = countD + countO + countSp
	rb := new(api.ResponseBody)
	rb.Result = countT
	api.WriteResponse(s, rb)
}
