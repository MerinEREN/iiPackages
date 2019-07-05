/*
Package photos "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package photos

import (
	"github.com/MerinEREN/iiPackages/datastore/photo"
	// "github.com/MerinEREN/iiPackages/storage"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler ...
func Handler(s *session.Session) {
	switch s.R.Method {
	/*
		case "POST":
			// https://stackoverflow.com/questions/15202448/go-formfile-for-multiple-files
			err := s.R.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile.
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			uID := s.R.Form.Get("uID")
			tagIDs := s.R.MultipartForm.Value["tagIDs"]
			description := s.R.MultipartForm.Value["description"][0]
			d := &photo.Photo{
				Description: description,
				Status:      "active",
			}
			fhx := s.R.MultipartForm.File["photos"]
			linksPhoto := make([]string, len(fhx))
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
				linksPhoto = append(linksPhoto, link)
			}
			d.LinksPhoto = linksPhoto
			pk, err := datastore.DecodeKey(uID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			k := datastore.NewIncompleteKey(s.Ctx, "Photo", pk)
			_, err = photo.Put(s, d, k)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", URL.Path, err)
				http.Error(s.W, err.Error(), http.StatusInternalServerError)
				return
			}
			s.W.WriteHeader(http.StatusNoContent)
	*/
	default:
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
		pType := q.Get("type")
		if pType == "" {
			/*
				ps, err := photo.GetFilteredByAncestor(s.Ctx, parentKey)
				if err != datastore.Done {
					log.Printf("Path: %s, Error: %v\n", URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				if len(ps) == 0 {
					s.W.WriteHeader(http.StatusNoContent)
					return
				}
			*/
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
