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
	"github.com/MerinEREN/iiPackages/datastore/servicePack"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/datastore/userTag"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
	"net/url"
)

// Handler returns account's servicePacks via account ID if provided.
// Otherwise if the user is "admin" returns servicePacks via account's all user's tag keys
// else returns only via user tag keys to show in timeline.
// If the request method is POST, uploads the files if present to the storage and
// puts the servicePack to the datastore.
func Handler(s *session.Session) {
	switch s.R.Method {
	case "POST":
		// https://stackoverflow.com/questions/15202448/go-formfile-for-multiple-files
		err := s.R.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile.
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		uID := s.R.Form.Get("uID")
		typ := s.R.MultipartForm.Value["type"][0]
		tagIDs := s.R.MultipartForm.Value["tagIDs"]
		title := s.R.MultipartForm.Value["title"][0]
		explanation := s.R.MultipartForm.Value["explanation"][0]
		sp := &servicePack.ServicePack{
			Type:        typ,
			TagIDs:      tagIDs,
			Title:       title,
			Explanation: explanation,
		}
		fhx := s.R.MultipartForm.File["photos"]
		linksPhoto := make([]string, len(fhx))
		for _, v := range fhx {
			f, err := v.Open()
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			defer f.Close()
			link, err := storage.UploadFile(s, f, v)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			linksPhoto = append(linksPhoto, link)
		}
		sp.LinksPhoto = linksPhoto
		err = servicePack.Put(s, sp, uID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
		return
	default:
		err := s.R.ParseForm()
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		accID := s.R.Form.Get("aID")
		var crsrAsStringx []string
		sps := make(servicePack.ServicePacks)
		URL := s.R.URL
		q := URL.Query()
		rb := new(api.ResponseBody)
		if accID != "" {
			aKey, err := datastore.DecodeKey(accID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			uKeyx, err := user.GetKeysByParentOrdered(s.Ctx, aKey)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			if len(s.R.Form["cs"]) == 0 {
				crsrAsStringx = make([]string, len(uKeyx))
			} else {
				crsrAsStringx = s.R.Form["cs"]
			}
			for i, v := range crsrAsStringx {
				sps2, crsrAsString, err := servicePack.GetByParent(s.Ctx, v, uKeyx[i])
				if err != nil && err != datastore.Done {
					log.Printf("Path: %s, Request: get account servicePacks via users keys, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i2, v2 := range sps2 {
					sps[i2] = v2
				}
				if i == 1 {
					q.Set("cs", crsrAsString)
				} else {
					q.Add("cs", crsrAsString)
				}
			}
			URL.RawQuery = q.Encode()
			rb.PrevPageURL = URL.String()
		} else {
			// For timeline
			var userTagKeys []*datastore.Key
			var accTagKeys []*datastore.Key
			var tagKeysQuery []*datastore.Key
			uKey, err := datastore.DecodeKey(s.R.Form.Get("uID"))
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
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
						s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
			} else {
				err = datastore.Get(s.Ctx, uKey, u)
				if err == datastore.ErrNoSuchEntity {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
					http.Error(s.W, err.Error(), http.StatusNoContent)
					return
				} else if err != nil {
					log.Printf("Path: %s, Error: %v\n",
						s.R.URL.Path, err)
					// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				} else {
					bs, err := json.Marshal(u)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "u",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
				}
			}
			isAdmin, err := u.IsAdmin(s.Ctx)
			if isAdmin {
				item, err = memcache.Get(s.Ctx, "accTagKeys")
				if err == nil {
					err = json.Unmarshal(item.Value, &accTagKeys)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					tagKeysQuery = accTagKeys
				} else {
					uKeys, err := user.GetKeysByParentOrdered(s.Ctx,
						uKey.Parent())
					if err == datastore.Done {
						if len(uKeys) == 0 {
							log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", s.R.URL.Path, err)
							s.W.WriteHeader(http.StatusNoContent)
							return
						}
					} else if err != nil {
						log.Printf("Path: %s, Request: get user keys via parent, Error: %v\n", s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for _, v := range uKeys {
						userTagKeys, err = userTag.GetKeysUserOrTag(s.Ctx, v)
						if err == datastore.Done {
							if len(userTagKeys) == 0 {
								log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
							}
						} else if err != nil {
							log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
							http.Error(s.W, err.Error(),
								http.StatusInternalServerError)
							return
						}
						for _, v2 := range userTagKeys {
							absent := true
							for _, v3 := range accTagKeys {
								if *v3 == *v2 {
									absent = false
								}
							}
							if absent {
								accTagKeys = append(accTagKeys, v2)
							}
						}
					}
					bs, err := json.Marshal(accTagKeys)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "accTagKeys",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
					tagKeysQuery = accTagKeys
				}
			} else {
				item, err = memcache.Get(s.Ctx, "userTagKeys")
				if err == nil {
					err = json.Unmarshal(item.Value, &tagKeysQuery)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
				} else {
					userTagKeys, err = userTag.GetKeysUserOrTag(s.Ctx, uKey)
					if err == datastore.Done {
						if len(userTagKeys) == 0 {
							log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
							s.W.WriteHeader(http.StatusNoContent)
							return
						}
					} else if err != nil {
						log.Printf("Path: %s, Request: getting user's tags, Error: %v\n", s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					bs, err := json.Marshal(userTagKeys)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
					item = &memcache.Item{
						Key:   "userTagKeys",
						Value: bs,
					}
					err = memcache.Add(s.Ctx, item)
					if err != nil {
						log.Printf("Path: %s, Error: %v\n",
							s.R.URL.Path, err)
					}
					tagKeysQuery = userTagKeys
				}
			}
			drct := s.R.Form.Get("d")
			if len(s.R.Form["cs"]) == 0 {
				crsrAsStringx = make([]string, len(tagKeysQuery))
			} else {
				crsrAsStringx = s.R.Form["cs"]
			}
			switch drct {
			case "next":
				for i, v := range crsrAsStringx {
					sps2, crsrAsString, err := servicePack.GetNewest(s.Ctx, v, tagKeysQuery[i])
					if err != nil && err != datastore.Done {
						log.Printf("Path: %s, Request: get newest servicePacks via user or users keys, Error: %v\n", s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range sps2 {
						sps[i2] = v2
					}
					if i == 1 {
						q.Set("cs", crsrAsString)
					} else {
						q.Add("cs", crsrAsString)
					}
				}
				URL.RawQuery = q.Encode()
				rb.NextPageURL = URL.String()
			case "prev":
				for i, v := range crsrAsStringx {
					sps2, crsrAsString, err := servicePack.GetOldest(s.Ctx, v, tagKeysQuery[i])
					if err != nil && err != datastore.Done {
						log.Printf("Path: %s, Request: get oldest servicePacks via user or users keys, Error: %v\n", s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range sps2 {
						sps[i2] = v2
					}
					if i == 1 {
						q.Set("cs", crsrAsString)
					} else {
						q.Add("cs", crsrAsString)
					}
				}
				URL.RawQuery = q.Encode()
				rb.PrevPageURL = URL.String()
			default:
				urlVsNext := url.Values{}
				urlVsPrev := url.Values{}
				for _, v := range tagKeysQuery {
					sps2, crsrNewAsString, crsrOldAsString, err := servicePack.Get(s.Ctx, v)
					if err != nil && err != datastore.Done {
						log.Printf("Path: %s, Request: get servicePacks via user or users keys, Error: %v\n", s.R.URL.Path, err)
						http.Error(s.W, err.Error(),
							http.StatusInternalServerError)
						return
					}
					for i2, v2 := range sps2 {
						sps[i2] = v2
					}
					urlVsNext.Add("cs", crsrNewAsString)
					urlVsPrev.Add("cs", crsrOldAsString)
				}
				urlVsNext.Set("uID", uKey.Encode())
				urlVsNext.Set("d", "next")
				urlVsPrev.Set("uID", uKey.Encode())
				urlVsPrev.Set("d", "prev")
				URL.RawQuery = urlVsPrev.Encode()
				rb.PrevPageURL = URL.String()
				URL.RawQuery = urlVsNext.Encode()
				rb.NextPageURL = URL.String()
				rb.Reset = true
			}
		}
		if len(sps) == 0 {
			s.W.WriteHeader(http.StatusNoContent)
			return
		}
		rb.Result = sps
		api.WriteResponseJSON(s, rb)
	}
}
