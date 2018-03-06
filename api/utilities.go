/*
Package api has utility functions to use with request handlers..
*/
package api

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/datastore/content"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
)

// WriteResponse "Exported functions should have a comment"
func WriteResponse(s *session.Session, rb *ResponseBody) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	s.W.Write(bs)
}

// GetLangValue sends the requested language value only.
func GetLangValue(cs content.Contents, code string) (content.Contents, error) {
	contentsClient := make(map[string]*content.Content)
	var err error
	for _, v := range cs {
		err = json.Unmarshal(v.ValuesBS, v.Values)
		if err != nil {
			return nil, err
		}
		contentsClient[v.ID].ID = v.ID
		contentsClient[v.ID].Value = v.Values[code]
		contentsClient[v.ID].LastModified = v.LastModified
		contentsClient[v.ID].Created = v.Created
	}
	return contentsClient, err
}
