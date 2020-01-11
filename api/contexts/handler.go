/*
Package contexts "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package contexts

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/context"
	// "github.com/MerinEREN/iiPackages/datastore/page"
	"github.com/MerinEREN/iiPackages/datastore/pageContext"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io"
	// "io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Handler handles contexts of pages and contexts page.
// ADD AUTHORISATION CONTROL !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	var err error
	switch s.R.Method {
	case "POST":
		var cx []*context.Context
		// Using 'decoder' is an alternative and can be used if response body has
		// more than one json object.
		// Otherwise don't use it, because it has performance disadvantages
		// compared to other solution.
		dec := json.NewDecoder(s.R.Body)
		for {
			c := new(context.Context)
			if err = dec.Decode(c); err == io.EOF {
				break
			} else if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			// map is not a valid datastore type to store
			// that is the reason of the marshalling below.
			c.ValuesBS, err = json.Marshal(c.Values)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			cx = append(cx, c)
		}
		// Get the entities from the begining with "" cursor string value.
		cs, crsrAsString, err := context.AddMultiAndGetNextLimited(s.Ctx, "", cx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		k := new(datastore.Key)
		for _, v := range cs {
			// If entity is from datastore.
			if len(v.PageIDs) == 0 {
				contextValues := make(map[string]string)
				err = json.Unmarshal(v.ValuesBS, &contextValues)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				v.Values = contextValues
				k, err = datastore.DecodeKey(v.ID)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				_, kpx, err := pageContext.GetKeysByPageOrContextKey(
					s.Ctx, k)
				if err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				var ekpx []string
				for _, v2 := range kpx {
					ekpx = append(ekpx, v2.Encode())
				}
				v.PageIDs = ekpx
			}
		}
		next := api.GenerateSubLink(s, crsrAsString, "next")
		s.W.Header().Set("Link", next)
		s.W.Header().Set("X-Reset", "true")
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponseJSON(s, cs)
	case "PUT":
		var cx []*context.Context
		dec := json.NewDecoder(s.R.Body)
		for {
			c := new(context.Context)
			if err = dec.Decode(c); err == io.EOF {
				break
			} else if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			// map is not a valid datastore type to store
			// that is the reason of the marshalling below.
			c.ValuesBS, err = json.Marshal(c.Values)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			cx = append(cx, c)
		}
		_, err = context.UpdateMulti(s.Ctx, cx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	case "DELETE":
		ekx := q["IDs"]
		if len(ekx) == 0 {
			log.Printf("Path: %s, Error: no key\n", URL.Path)
			http.Error(s.W, "no key", http.StatusBadRequest)
			return
		}
		err = context.DeleteMulti(s.Ctx, ekx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		cs := make(context.Contexts)
		if pID := q.Get("pID"); pID != "" {
			kp := datastore.NewKey(s.Ctx, "Page", pID, 0, nil)
			_, kx, err := pageContext.GetKeysByPageOrContextKey(s.Ctx, kp)
			if err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			if kx != nil {
				var cx []*context.Context
				// RETURNED ENTITY LIMIT COULD BE A PROBLEM HERE !!!!!!!!!!
				err = datastore.GetMulti(s.Ctx, kx, cx)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path,
						err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i, v := range kx {
					cs[v.Encode()] = cx[i]
				}
				cookie, err := s.R.Cookie("lang")
				if err == http.ErrNoCookie {
					// "http.Status..." MAY BE WRONG !!!!!!!!!!!!!!!!!!
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusBadRequest)
					return
				}
				rb, err := GetLangValue(cs, cookie.Value)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				s.W.Header().Set("X-Reset", "true")
				s.W.Header().Set("Content-Type", "application/json")
				api.WriteResponseJSON(s, rb)
			} else {
				s.W.WriteHeader(http.StatusNoContent)
			}
		} else {
			limit := q.Get("limit")
			var lim int
			if limit == "" {
				// lim = 0
				lim = 2222
			} else {
				lim, err = strconv.Atoi(limit)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						URL.Path, err)
				}
			}
			crsrAsString := q.Get("after")
			cs, crsrAsString, err = context.GetNextLimited(s.Ctx,
				crsrAsString, lim)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			if len(cs) == 0 {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
			k := new(datastore.Key)
			for _, v := range cs {
				contextValues := make(map[string]string)
				err = json.Unmarshal(v.ValuesBS, &contextValues)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				v.Values = contextValues
				k, err = datastore.DecodeKey(v.ID)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				_, kpx, err := pageContext.GetKeysByPageOrContextKey(
					s.Ctx, k)
				if err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for _, v2 := range kpx {
					v.PageIDs = append(v.PageIDs, v2.Encode())
				}
			}
			next := api.GenerateSubLink(s, crsrAsString, "next")
			s.W.Header().Set("Link", next)
			s.W.Header().Set("Content-Type", "application/json")
			api.WriteResponseJSON(s, cs)
		}
	}
}
