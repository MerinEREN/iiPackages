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

// WriteResponse writes a JSON encoded http response
// with "text/plain" value as the default "Content-Type".
func WriteResponse(s *session.Session, rb *ResponseBody) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	s.W.Write(bs)
}

// WriteResponseJSON writes a JSON encoded http response
// with "application/json" value as "Content-Type".
func WriteResponseJSON(s *session.Session, rb *ResponseBody) {
	s.W.Header().Set("Content-Type", "application/json")
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	s.W.Write(bs)
}

// GetLangValue sends the requested language value of the contents only.
func GetLangValue(cs content.Contents, lang string) (map[string]string, error) {
	contentsClient := make(map[string]string)
	var err error
	for i, v := range cs {
		contentValues := make(map[string]string)
		err = json.Unmarshal(v.ValuesBS, &contentValues)
		if err != nil {
			return nil, err
		}
		contentsClient[i] = contentValues[lang]
	}
	return contentsClient, err
}
