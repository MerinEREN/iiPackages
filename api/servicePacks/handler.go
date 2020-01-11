/*
Package servicePacks "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package servicePacks

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/servicePack"
	"github.com/MerinEREN/iiPackages/datastore/tagServicePack"
	"github.com/MerinEREN/iiPackages/datastore/tagUser"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handler returns account's servicePacks via account ID if provided.
// Otherwise if the user is "admin" returns servicePacks via account's all user's tag keys
// else returns only via user tag keys to show in timeline.
// If the request method is POST, uploads the files if present to the storage and
// puts the servicePack to the datastore.
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
		typ := s.R.MultipartForm.Value["type"][0]
		tIDs := s.R.MultipartForm.Value["tagIDs"]
		title := s.R.MultipartForm.Value["title"][0]
		description := s.R.MultipartForm.Value["description"][0]
		sp := &servicePack.ServicePack{
			Type:         typ,
			Title:        title,
			Description:  description,
			Status:       "active",
			Created:      time.Now(),
			LastModified: time.Now(),
		}
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
		pk, err := datastore.DecodeKey(uID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		k := datastore.NewIncompleteKey(s.Ctx, "ServicePack", pk)
		tspx := make([]*tagServicePack.TagServicePack, 0, cap(tIDs))
		ktspx := make([]*datastore.Key, 0, cap(tIDs))
		kpx := make([]*datastore.Key, 0, cap(px))
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			k, err1 = datastore.Put(ctx, k, sp)
			if err1 != nil {
				return
			}
			for _, v := range tIDs {
				ktsp := datastore.NewKey(s.Ctx, "TagServicePack", v, 0, k)
				ktspx = append(ktspx, ktsp)
				kt := new(datastore.Key)
				kt, err1 = datastore.DecodeKey(v)
				if err1 != nil {
					return
				}
				tsp := &tagServicePack.TagServicePack{
					Created: time.Now(),
					TagKey:  kt,
				}
				tspx = append(tspx, tsp)
			}
			_, err1 = datastore.PutMulti(ctx, ktspx, tspx)
			if err1 != nil {
				return
			}
			for i := 0; i < len(px); i++ {
				kp := datastore.NewIncompleteKey(s.Ctx, "Photo", k)
				kpx = append(kpx, kp)
			}
			_, err1 = datastore.PutMulti(ctx, kpx, px)
			return
		}, nil)
		if err != nil {
			// REMOVE ALL THE UPLOADED FILES FROM THE STORAGE !!!!!!!!!!!!!!!!!
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		accID := q.Get("aID")
		var crsrAsStringx []string
		sps := make(servicePack.ServicePacks)
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
				sps2, crsrAsString, err := servicePack.GetNextByParentLimited(s.Ctx, after[i], v, lim)
				if err != nil {
					log.Printf("Path: %s, Request: get account servicePacks via users keys, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i2, v2 := range sps2 {
					sps[i2] = v2
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
					kx, crsrAsString, err := tagServicePack.GetPrevKeysParentsFilteredByTagKey(s.Ctx, before[i], v)
					if err != datastore.Done {
						log.Printf("Path: %s, Request: get previous service pack's keys by tag key, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					sps2, err := servicePack.GetMulti(s.Ctx, kx)
					if err != nil {
						log.Printf("Path: %s, Request: get previous service packs, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range sps2 {
						sps[i2] = v2
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
					kx, crsrAsString, err := tagServicePack.GetNextKeysParentsFilteredByTagKeyLimited(s.Ctx, after[i], v, lim)
					if err != nil {
						log.Printf("Path: %s, Request: get next service pack's keys by tag key, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					sps2, err := servicePack.GetMulti(s.Ctx, kx)
					if err != nil {
						log.Printf("Path: %s, Request: get next service packs, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range sps2 {
						sps[i2] = v2
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
					kx, beforeAsString, afterAsString, err := tagServicePack.GetKeysParentsFilteredByTagKeyLimited(s.Ctx, v, lim)
					if err != nil {
						log.Printf("Path: %s, Request: get initial service pack's keys by tag key, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					sps2, err := servicePack.GetMulti(s.Ctx, kx)
					if err != nil {
						log.Printf("Path: %s, Request: get initial service packs, Error: %v\n", URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range sps2 {
						sps[i2] = v2
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
		if len(sps) == 0 {
			s.W.WriteHeader(http.StatusNoContent)
			return
		}
		s.W.Header().Set("Content-Type", "application/json")
		api.WriteResponseJSON(s, sps)
	}
}
