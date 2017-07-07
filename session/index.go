/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package session

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"net/http"
)

func (s *Session) Init(w http.ResponseWriter, r *http.Request) {
	s.Ctx = appengine.NewContext(r)
	s.R = r
	s.W = w
	s.U = user.Current(s.Ctx)
}
