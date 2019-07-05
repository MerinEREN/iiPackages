/*
Package api has utility functions to use with request handlers..
*/
package api

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// GenerateSubLink generates next and prev links for response's Link header.
func GenerateSubLink(s *session.Session, i interface{}, rel string) string {
	URL := s.R.URL
	q := URL.Query()
	switch rel {
	case "prev":
		switch v := i.(type) {
		case string:
			q.Set("before", v)
		case []string:
			for i, v2 := range v {
				if i == 0 {
					q.Set("before", v2)
				} else {
					q.Add("before", v2)
				}
			}
		default:
			// if i is datastore.Cursor type
			crsr := v.(datastore.Cursor)
			q.Set("before", crsr.String())
		}
	default:
		// next
		switch v := i.(type) {
		case string:
			q.Set("after", v)
		case []string:
			for i, v2 := range v {
				if i == 0 {
					q.Set("after", v2)
				} else {
					q.Add("after", v2)
				}
			}
		default:
			// if i is datastore.Cursor type
			crsr := v.(datastore.Cursor)
			q.Set("after", crsr.String())
		}
	}
	URL.RawQuery = q.Encode()
	pageURIRef := "<" + URL.String() + ">"
	sx := []string{"rel", rel}
	rel = strings.Join(sx, "=")
	sx[0] = pageURIRef
	sx[1] = rel
	return strings.Join(sx, "; ")
}

// WriteResponseJSON writes a JSON encoded http response.
func WriteResponseJSON(s *session.Session, rb interface{}) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	s.W.Write(bs)
}
