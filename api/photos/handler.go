/*
Package photos "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package photos

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// Handler ...
func Handler(s *session.Session) {
	URL := s.R.URL
	q := URL.Query()
	pID := q.Get("pID")
	if pID == "" {
		log.Printf("Path: %s, Error: no parent ID\n", URL.Path)
		http.Error(s.W, "No parent ID", http.StatusBadRequest)
		return
	}
	pk, err := datastore.DecodeKey(pID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	switch s.R.Method {
	case "POST":
		ct := s.R.Header.Get("Content-Type")
		if ct != "application/json" {
			log.Printf("Path: %s, Error: %v\n", URL.Path, "Content type is not application/json")
			http.Error(s.W, "Content type is not application/json",
				http.StatusUnsupportedMediaType)
			return
		}
		var bs []byte
		bs, err = ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		var fhx []*multipart.FileHeader
		err = json.Unmarshal(bs, &fhx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		px := make([]*photo.Photo, 0, cap(fhx))
		kx := make([]*datastore.Key, 0, cap(fhx))
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
				Link:         link,
				Status:       "pending",
				Created:      time.Now(),
				LastModified: time.Now(),
			}
			px = append(px, p)
			k := datastore.NewIncompleteKey(s.Ctx, "Photo", pk)
			kx = append(kx, k)
		}
		switch len(kx) {
		case 1:
			_, err = datastore.Put(s.Ctx, kx[0], px[0])
			if err != nil {
				// REMOVE THE UPLOADED FILE FROM THE STORAGE !!!!!!!!!!!!!!
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
		default:
			_, err = datastore.PutMulti(s.Ctx, kx, px)
			if err != nil {
				// REMOVE ALL THE UPLOADED FILES FROM THE STORAGE !!!!!!!!!
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
		}
		/*
			NOT RETURNING PHOTOS, BECAUSE AT THE FINAL VERSION
			THE INITIAL VALUE OF THE "status" OF THE PHOTO WON'T BE "active"
			AND PHOTO'S STATUS WILL BE "active" AFTER CONFIRMATION PROCESS.
		*/
		// SEND THE INFORMATION TO SHOWN AS SNACKBAR MESSAGE !!!!!!!!!!!!!!!!!!!!!!
		s.W.WriteHeader(http.StatusNoContent)
	default:
		pType := q.Get("type")
		if pType == "" {
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
			ps, err := photo.GetFilteredByAncestorLimited(s.Ctx, pk, lim)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			switch len(ps) {
			case 0:
				s.W.WriteHeader(http.StatusNoContent)
			case 1:
				s.W.Header().Set("Content-Type", "application/json")
				for _, v := range ps {
					api.WriteResponseJSON(s, v)
				}
			default:
				s.W.Header().Set("Content-Type", "application/json")
				api.WriteResponseJSON(s, ps)
			}
		} else {
			p, err := photo.GetMainByAncestor(s.Ctx, pType, pk)
			if err == datastore.Done {
				s.W.WriteHeader(http.StatusNoContent)
				return
			}
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			s.W.Header().Set("Content-Type", "application/json")
			api.WriteResponseJSON(s, p)
		}
	}
}
