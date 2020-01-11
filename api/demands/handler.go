/*
Package demands "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package demands

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/tagDemand"
	"github.com/MerinEREN/iiPackages/datastore/tagUser"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
Handler returns account's demands via account ID if provided.
Otherwise if the user is "admin" returns demands via account's all user's tag keys
else returns only via user tag keys to show in timeline.
If the request method is POST, uploads the files if present to the storage and
puts the demand and it's tags to the datastore.
Also returns created demand.
*/
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	switch s.R.Method {
	case "POST":
		uID := q.Get("uID")
		if uID == "" {
			log.Printf("Path: %s, Error: No user ID.\n", URL.Path)
			http.Error(s.W, "No user ID.", http.StatusBadRequest)
			return
		}
		// https://stackoverflow.com/questions/15202448/go-formfile-for-multiple-files
		err := s.R.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile.
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		description := s.R.MultipartForm.Value["description"][0]
		d := &demand.Demand{
			Description:  description,
			Status:       "pending",
			Created:      time.Now(),
			LastModified: time.Now(),
		}
		tIDs := s.R.MultipartForm.Value["tagIDs"]
		fhx := s.R.MultipartForm.File["files"]
		px := make([]*photo.Photo, 0, cap(fhx))
		for _, v := range fhx {
			f, err := v.Open()
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			defer f.Close()
			link, err := storage.UploadFile(s, f, v)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			p := &photo.Photo{
				Link:   link,
				Status: "active",
			}
			px = append(px, p)
		}
		ku, err := datastore.DecodeKey(uID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		k := datastore.NewIncompleteKey(s.Ctx, "Demand", ku)
		k, err = demand.Add(s.Ctx, k, d, tIDs, px)
		if err != nil {
			// REMOVE ALL THE UPLOADED FILES FROM THE STORAGE !!!!!!!!!!!!!!!!!
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		d.ID = k.Encode()
		d.UserID = uID
		ka := ku.Parent()
		d.AccountID = ka.Encode()
		sx := []string{"/demands", d.ID}
		path := strings.Join(sx, "/")
		rel, err := URL.Parse(path)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.Header().Set("Location", rel.String())
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponseJSON(s, d)
	default:
		accID := q.Get("aID")
		var crsrAsStringx []string
		ds := make(demand.Demands)
		if accID != "" {
			ka, err := datastore.DecodeKey(accID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			kux, err := user.GetKeysByParentOrdered(s.Ctx, ka)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			after := q["after"]
			if len(after) == 0 {
				after = make([]string, len(kux))
			}
			var lim int
			limit := q.Get("limit")
			if limit == "" {
				lim = 0
			} else {
				lim, err = strconv.Atoi(limit)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				}
			}
			for i, v := range kux {
				ds2, crsrAsString, err := demand.GetNextByParentLimited(
					s.Ctx, after[i], v, lim)
				if err != nil {
					log.Printf("Path: %s, Request: get account demands via users keys, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i2, v2 := range ds2 {
					ds[i2] = v2
				}
				crsrAsStringx = append(crsrAsStringx, crsrAsString)
			}
			next := api.GenerateSubLink(s, crsrAsStringx, "next")
			s.W.Header().Set("Link", next)
		} else {
			// For timeline
			uID := q.Get("uID")
			if uID == "" {
				log.Printf("Path: %s, Error: No user ID.\n", URL.Path)
				http.Error(s.W, "No user ID.", http.StatusBadRequest)
				return
			}
			ku, err := datastore.DecodeKey(uID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			u := new(user.User)
			item, err := memcache.Get(s.Ctx, "u")
			if err == nil {
				err = json.Unmarshal(item.Value, u)
				if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				err = datastore.Get(s.Ctx, ku, u)
				if err == datastore.ErrNoSuchEntity {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
					http.Error(s.W, err.Error(), http.StatusNoContent)
					return
				} else if err != nil {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				} else {
					bs, err := json.Marshal(u)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "u",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
				}
			}
			var ktux []*datastore.Key
			var ktax []*datastore.Key
			var ktx []*datastore.Key
			isAdmin, err := u.IsAdmin(s.Ctx)
			if isAdmin {
				item, err = memcache.Get(s.Ctx, "ktax")
				if err == nil {
					err = json.Unmarshal(item.Value, &ktax)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					ktx = ktax
				} else {
					kux, err := user.GetKeysByParentOrdered(s.Ctx,
						ku.Parent())
					if err != nil {
						log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					if len(kux) == 0 {
						// Impossible !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
						log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", URL.Path, err)
						s.W.WriteHeader(http.StatusNoContent)
						return
					}
					for _, v := range kux {
						ktux, err = tagUser.GetKeysByUserOrTagKey(s.Ctx, v)
						if err == datastore.Done {
							if len(ktux) == 0 {
								log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
							}
						} else if err != nil {
							log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
							http.Error(s.W, err.Error(),
								http.StatusInternalServerError)
							return
						}
						for _, v2 := range ktux {
							absent := true
							for _, v3 := range ktax {
								if *v3 == *v2 {
									absent = false
								}
							}
							if absent {
								ktax = append(ktax, v2)
							}
						}
					}
					bs, err := json.Marshal(ktax)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "ktax",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
					ktx = ktax
				}
			} else {
				item, err = memcache.Get(s.Ctx, "ktux")
				if err == nil {
					err = json.Unmarshal(item.Value, &ktx)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
				} else {
					ktux, err = tagUser.GetKeysByUserOrTagKey(s.Ctx, ku)
					if err == datastore.Done {
						if len(ktux) == 0 {
							log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
							s.W.WriteHeader(http.StatusNoContent)
							return
						}
					} else if err != nil {
						log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					bs, err := json.Marshal(ktux)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "ktux",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
					ktx = ktux
				}
			}
			before := q["before"]
			after := q["after"]
			var crsrAsStringx []string
			if len(before) != 0 {
				for i, v := range ktx {
					kx, crsrAsString, err := tagDemand.GetPrevKeysParentsFilteredByTagKey(s.Ctx, before[i], v)
					if err != datastore.Done {
						log.Printf("Path: %s, Request: get previous demand's keys by tag key, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					ds2, err := demand.GetMulti(s.Ctx, kx)
					if err != nil {
						log.Printf("Path: %s, Request: get previous demands, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range ds2 {
						ds[i2] = v2
					}
					crsrAsStringx = append(crsrAsStringx, crsrAsString)
				}
				prev := api.GenerateSubLink(s, crsrAsStringx, "prev")
				s.W.Header().Set("Link", prev)
			} else if len(after) != 0 {
				limit := q.Get("limit")
				var lim int
				if limit == "" {
					lim = 0
				} else {
					lim, err = strconv.Atoi(limit)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
				}
				for i, v := range ktx {
					kx, crsrAsString, err := tagDemand.GetNextKeysParentsFilteredByTagKeyLimited(s.Ctx, after[i], v, lim)
					if err != nil {
						log.Printf("Path: %s, Request: get next demand's keys by tag key, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					ds2, err := demand.GetMulti(s.Ctx, kx)
					if err != nil {
						log.Printf("Path: %s, Request: get next demands, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range ds2 {
						ds[i2] = v2
					}
					crsrAsStringx = append(crsrAsStringx, crsrAsString)
				}
				next := api.GenerateSubLink(s, crsrAsStringx, "next")
				s.W.Header().Set("Link", next)
			} else {
				limit := q.Get("limit")
				var lim int
				if limit == "" {
					lim = 0
				} else {
					lim, err = strconv.Atoi(limit)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							URL.Path, err)
					}
				}
				var crsrAsStringx2 []string
				for _, v := range ktx {
					kx, beforeAsString, afterAsString, err := tagDemand.GetKeysParentsFilteredByTagKeyLimited(s.Ctx, v, lim)
					if err != nil {
						log.Printf("Path: %s, Request: get initial demand's keys by tag key, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					ds2, err := demand.GetMulti(s.Ctx, kx)
					if err != nil {
						log.Printf("Path: %s, Request: get initial demands, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range ds2 {
						ds[i2] = v2
					}
					crsrAsStringx = append(crsrAsStringx, beforeAsString)
					crsrAsStringx2 = append(crsrAsStringx2, afterAsString)
				}
				prev := api.GenerateSubLink(s, crsrAsStringx, "prev")
				next := api.GenerateSubLink(s, crsrAsStringx2, "next")
				sx := []string{prev, next}
				link := strings.Join(sx, ", ")
				s.W.Header().Set("Link", link)
				s.W.Header().Set("X-Reset", "true")
			}
		}
		switch len(ds) {
		case 0:
			s.W.WriteHeader(http.StatusNoContent)
		case 1:
			s.W.Header().Set("Content-Type", "application/json")
			for _, v := range ds {
				api.WriteResponseJSON(s, v)
			}
		default:
			s.W.Header().Set("Content-Type", "application/json")
			api.WriteResponseJSON(s, ds)
		}
	}
}
