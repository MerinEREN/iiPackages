// Package tags posts a tag and also gets all tags.
package tags

import (
	// "github.com/MerinEREN/iiPackages/api/user"
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/content"
	"github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"sort"
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
	isContentEditor, err := u.IsContentEditor(s.Ctx)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Type != "inHouse" || !(isAdmin || isContentEditor) {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are unauthorized user.", http.StatusUnauthorized)
		return
	} */
	ts := make(tag.Tags)
	var err error
	switch s.R.Method {
	case "POST":
		contentID := s.R.FormValue("contentID")
		t := &tag.Tag{
			ContentID: contentID,
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
		st := s.R.FormValue("st")
		if st != "" {
			var kx []*datastore.Key
			count := 0
			switch st {
			case "top":
				// Get most used by demands
				dx, err := demand.GetAllLimited(s.Ctx, 100)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				if len(dx) == 0 {
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
				var elem int
				var ok bool
				tagIDsMap := make(map[string]int)
				for _, v := range dx {
					for _, v2 := range v.TagIDs {
						if elem, ok = tagIDsMap[v2]; ok {
							tagIDsMap[v2] = elem + 1
						} else {
							tagIDsMap[v2] = 1
						}
					}
				}
				var twcx []TagWithCount
				for i, v := range tagIDsMap {
					twc := TagWithCount{i, v}
					twcx = append(twcx, twc)
				}
				sort.Slice(twcx, func(i, j int) bool {
					return twcx[i].Count > twcx[j].Count
				})
				for _, v := range twcx {
					if count == 6 {
						break
					}
					k, err := datastore.DecodeKey(v.TagID)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					kx = append(kx, k)
					count = count + 1
				}
			default:
				// Get filtered via search text(st)
				// Get content keys from all tag keys.
				kcx, err := tag.GetAllDecodedStringIDs(s.Ctx)
				if err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
						err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				cx := make([]content.Content, len(kcx))
				err = datastore.GetMulti(s.Ctx, kcx, cx)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
						err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				cookie, err := s.R.Cookie("lang")
				if err == http.ErrNoCookie {
					log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
						err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i, v := range cx {
					if count == 6 {
						break
					}
					contentValues := make(map[string]string)
					err = json.Unmarshal(v.ValuesBS, &contentValues)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return

					}
					v.Value = contentValues[cookie.Value]
					if strings.Contains(strings.ToLower(v.Value),
						strings.ToLower(st)) {
						k := datastore.NewKey(s.Ctx, "Tag",
							kcx[i].Encode(), 0, nil)
						kx = append(kx, k)
						count = count + 1
					}
				}
			}
			for _, v := range kx {
				t := new(tag.Tag)
				t.ContentID = v.StringID()
				t.ID = v.Encode()
				ts[t.ID] = t
			}
		} else {
			// Get all
			ts, err = tag.GetMulti(s.Ctx, nil)
			if err != datastore.Done {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
		}
		if len(ts) == 0 {
			s.W.WriteHeader(http.StatusNoContent)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = ts
		api.WriteResponseJSON(s, rb)
	}
}
