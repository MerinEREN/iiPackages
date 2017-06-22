package apis

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
)

func WriteResponse(s *session.Session, rb *ResponseBody) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	s.W.Write(bs)
}
