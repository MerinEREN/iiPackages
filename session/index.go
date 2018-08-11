package session

import (
	"google.golang.org/appengine"
	userG "google.golang.org/appengine/user"
	"net/http"
)

// Init initialises the Session struct.
func (s *Session) Init(w http.ResponseWriter, r *http.Request) {
	s.Ctx = appengine.NewContext(r)
	s.W = w
	s.R = r
	s.U = userG.Current(s.Ctx)
}
