// Package tags posts a tag and also gets all tags.
package tags

import (
	// "github.com/MerinEREN/iiPackages/api/user"
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/context"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/datastore/tagDemand"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// TagWithCount is a struct which containes most recent demand's used tag's ID with their
// usage counts.
type TagWithCount struct {
	TagID string
	Count int
}

// Handler posts a tag and returns all the tags from the begining of the kind
// and gets all the tags from the begining of the kind.
func Handler(s *session.Session) {
	// THE CONTROLS BELOVE PREVENT GET REQUEST THAT NECESSARY FOR SELECT SOME SELECT
	// FIELDS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* u, err := user.Get(s)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status == "suspended" {
		log.Printf("Suspended user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are suspended", http.StatusForbidden)
		return
	}
	isAdmin, err := u.IsAdmin(s.Ctx)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	isContextEditor, err := u.IsContextEditor(s.Ctx)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Type != "inHouse" || !(isAdmin || isContextEditor) {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are unauthorized user.", http.StatusUnauthorized)
		return
	} */
	ts := make(tag.Tags)
	var err error
	switch s.R.Method {
	case "POST":
		contextID := s.R.FormValue("contextID")
		t := &tag.Tag{
			ContextID: contextID,
		}
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		ts, err = tag.PutAndGetMulti(s, t)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = ts
		rb.PrevPageURL = "/tags?c=" + crsr.String()
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	default:
		// Handles "GET" requests
		URL := s.R.URL
		qry := URL.Query()
		q := qry.Get("q")
		if q != "" {
			var kx []*datastore.Key
			count := 0
			switch q {
			case "top":
				// Get most used by demands
				tdx, err := tagDemand.GetDistinctLatestLimited(s.Ctx, 100)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				if len(tdx) == 0 {
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
				for _, v := range tdx {
					kx = append(kx, v.TagKey)
				}
			default:
				// Get filtered via search text(q)
				k2x, err := tag.GetAllKeysOnly(s.Ctx)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				var kcx []*datastore.Key
				for _, v := range k2x {
					k, err := datastore.DecodeKey(v.StringID())
					if err != nil {
						log.Printf("Path: %s, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					kcx = append(kcx, k)
				}
				cx := make([]context.Context, len(kcx))
				err = datastore.GetMulti(s.Ctx, kcx, cx)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				cookie, err := s.R.Cookie("lang")
				if err == http.ErrNoCookie {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i, v := range cx {
					if count == 6 {
						break
					}
					contextValues := make(map[string]string)
					err = json.Unmarshal(v.ValuesBS, &contextValues)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return

					}
					v.Value = contextValues[cookie.Value]
					if strings.Contains(strings.ToLower(v.Value),
						strings.ToLower(q)) {
						k := datastore.NewKey(s.Ctx, "Tag",
							kcx[i].Encode(), 0, nil)
						kx = append(kx, k)
						count = count + 1
					}
				}
			}
			for _, v := range kx {
				t := new(tag.Tag)
				t.ContextID = v.StringID()
				t.ID = v.Encode()
				ts[t.ID] = t
			}
		} else {
			// Get all
			ts, err = tag.GetMulti(s.Ctx, nil)
			if err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
		}
		if len(ts) == 0 {
			s.W.WriteHeader(http.StatusNoContent)
			return
		}
		s.W.Header().Set("Content-Type", "application/json")
		api.WriteResponseJSON(s, ts)
	}
}
