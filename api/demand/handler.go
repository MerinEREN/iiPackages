// Package demand updates and returns a demand.
package demand

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/storage"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// Handler updates and returns demand via demand ID which is an encoded key.
func Handler(s *session.Session) {
	rb := new(api.ResponseBody)
	ID := strings.Split(s.R.URL.Path, "/")[2]
	if ID == "" {
		log.Printf("Path: %s, Error: no demand ID\n", s.R.URL.Path)
		http.Error(s.W, "No demand ID", http.StatusBadRequest)
		return
	}
	d := new(demand.Demand)
	k := new(datastore.Key)
	var err error
	switch s.R.Method {
	case "DELETE":
		err = demand.UpdateStatus(s.Ctx, ID, "deleted")
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
		return
	case "PUT":
		// IF "Content-Type" HEADER IS NOT "application/json" THROW A
		// "415 Unsupported Media Type HTTP status code" !!!!!!!!!!!!!!!!!!!!!!!!!!
		// https://stackoverflow.com/questions/15202448/go-formfile-for-multiple-files
		err = s.R.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile.
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		tagIDs := s.R.MultipartForm.Value["tagIDs"]
		description := s.R.MultipartForm.Value["description"][0]
		d = &demand.Demand{
			Description: description,
			Status:      "modified",
		}
		// DELETE UNUSED FILES IN STORAGE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
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
		d.LinksPhoto = linksPhoto
		k, err = datastore.DecodeKey(ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		d, err = demand.Put(s, d, k)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		k, err = datastore.DecodeKey(ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		err = datastore.Get(s.Ctx, k, d)
		if err == datastore.ErrNoSuchEntity {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(),
				http.StatusInternalServerError)
			return
		}
	}
	ds := make(demand.Demands)
	d.ID = ID
	d.UserID = k.Parent().Encode()
	d.AccountID = k.Parent().Parent().Encode()
	ds[ID] = d
	rb.Result = ds
	api.WriteResponseJSON(s, rb)
}
