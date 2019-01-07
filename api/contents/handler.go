/*
Package contents "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package contents

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/content"
	// "github.com/MerinEREN/iiPackages/datastore/page"
	"github.com/MerinEREN/iiPackages/datastore/pageContent"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"io"
	// "io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Handler handles contents of pages and contents page.
// ADD AUTHORISATION CONTROL !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Handler(s *session.Session) {
	var err error
	switch s.R.Method {
	case "POST":
		var cx []*content.Content
		// Using 'decoder' is an alternative and can be used if response body has
		// more than one json object.
		// Otherwise don't use it, because it has performance disadvantages
		// compared to other solution.
		dec := json.NewDecoder(s.R.Body)
		for {
			c := new(content.Content)
			if err = dec.Decode(c); err == io.EOF {
				break
			} else if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			// map is not a valid datastore type to store
			// that is the reason of the marshalling below.
			c.ValuesBS, err = json.Marshal(c.Values)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			cx = append(cx, c)
		}
		cs := make(content.Contents)
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		cs, crsr, err = content.PutMultiAndGetMulti(s, crsr, cx)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		k := new(datastore.Key)
		for _, v := range cs {
			// If entity is from datastore.
			if len(v.PageIDs) == 0 {
				contentValues := make(map[string]string)
				err = json.Unmarshal(v.ValuesBS, &contentValues)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Values = contentValues
				k, err = datastore.DecodeKey(v.ID)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				_, kpx, err := pageContent.GetKeysWithPageOrContentKeys(s.Ctx, k)
				if err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				var ekpx []string
				for _, v2 := range kpx {
					ekpx = append(ekpx, v2.Encode())
				}
				v.PageIDs = ekpx
			}
		}
		rb := new(api.ResponseBody)
		rb.Result = cs
		rb.PrevPageURL = "/contents?c=" + crsr.String()
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	case "PUT":
		var cx []*content.Content
		dec := json.NewDecoder(s.R.Body)
		for {
			c := new(content.Content)
			if err = dec.Decode(c); err == io.EOF {
				break
			} else if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			// map is not a valid datastore type to store
			// that is the reason of the marshalling below.
			c.ValuesBS, err = json.Marshal(c.Values)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			cx = append(cx, c)
		}
		_, err = content.PutMulti(s, cx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	case "DELETE":
		IDsAsString := s.R.FormValue("IDs")
		ekx := strings.Split(IDsAsString, ",")
		if len(ekx) == 0 {
			log.Printf("Path: %s, Error: no key\n", s.R.URL.Path)
			http.Error(s.W, "no key", http.StatusBadRequest)
			return
		}
		err = content.DeleteMulti(s.Ctx, ekx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		crsr, err := datastore.DecodeCursor(s.R.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie, err := s.R.Cookie("lang")
		if err == http.ErrNoCookie {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		cs := make(content.Contents)
		rb := new(api.ResponseBody)
		if pID := s.R.FormValue("pageID"); pID != "" {
			k := datastore.NewKey(s.Ctx, "Page", pID, 0, nil)
			_, kx, err := pageContent.GetKeysWithPageOrContentKeys(s.Ctx, k)
			if err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			if kx != nil {
				cs, _, err = content.GetMulti(s, crsr, nil, kx)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				contentsClient, err := GetLangValue(cs, cookie.Value)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				rb.Result = contentsClient
			} else {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
		} else {
			cs, crsr, err = content.GetMulti(s, crsr, nil, nil)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(cs) == 0 {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
			k := new(datastore.Key)
			for _, v := range cs {
				contentValues := make(map[string]string)
				err = json.Unmarshal(v.ValuesBS, &contentValues)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Values = contentValues
				k, err = datastore.DecodeKey(v.ID)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				_, kpx, err := pageContent.GetKeysWithPageOrContentKeys(s.Ctx, k)
				if err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				for _, v2 := range kpx {
					v.PageIDs = append(v.PageIDs, v2.Encode())
				}
			}
			rb.PrevPageURL = "/contents?c=" + crsr.String()
			rb.Result = cs
		}
		api.WriteResponseJSON(s, rb)
	}
}
