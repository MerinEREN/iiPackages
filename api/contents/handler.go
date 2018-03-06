/*
Package contents "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
// GET LANGUAGE CODE AND PAGE NAME FROM URL AND HANDLE THE REQUEST HERE.
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
	"strconv"
)

// Handler handles contents of pages and contents page.
func Handler(s *session.Session) {
	var err error
	switch s.R.Method {
	case "POST":
		rb := new(api.ResponseBody)
		cs := make(content.Contents)
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
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		cs, crsr, err = content.PutMultiAndGetMulti(s, crsr, cx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v := range cs {
			// If entity is from datastore.
			if len(v.Pages) == 0 {
				contentValues := make(map[string]string)
				err = json.Unmarshal(v.ValuesBS, &contentValues)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Values = contentValues
				IDx, err := pageContent.Get(s, "", v.ID)
				if err != nil && err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Pages = IDx
			}
		}
		rb.Result = cs
		rb.Reset = true
		rb.PrevPageURL = "/pages?d=prev&" + "c=" + crsr.String()
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
		s.W.WriteHeader(http.StatusOK)
	default:
		rb := new(api.ResponseBody)
		cs := make(content.Contents)
		crsr, err := datastore.DecodeCursor(s.R.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		if pID := s.R.FormValue("pageID"); pID != "" {
			IDx, err := pageContent.Get(s, pID, "")
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			var kx []*datastore.Key
			var contentKeyIntID int64
			for _, v := range IDx {
				contentKeyIntID, err = strconv.ParseInt(v, 11, 64)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				kx = append(kx, datastore.NewKey(s.Ctx, "Content", "", contentKeyIntID, nil))
			}
			if kx != nil {
				cs, _, err = content.GetMulti(s, crsr, kx, nil)
				if err != nil && err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				lCode := s.R.FormValue("lCode")
				contentsClient, err := api.GetLangValue(cs, lCode)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				rb.Result = contentsClient
			} else {
				rb.Result = nil
			}
		} else {
			cs, crsr, err = content.GetMulti(s, crsr, nil, nil)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, v := range cs {
				contentValues := make(map[string]string)
				err = json.Unmarshal(v.ValuesBS, &contentValues)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Values = contentValues
				IDx, err := pageContent.Get(s, "", v.ID)
				if err != nil && err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Pages = IDx
			}
			rb.PrevPageURL = "/contents?c=" + crsr.String()
			rb.Result = cs
		}
		api.WriteResponse(s, rb)
	}
}
