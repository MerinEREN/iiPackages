/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package session

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
	"net/http"
)

type Session struct {
	Ctx context.Context
	R   *http.Request
	W   http.ResponseWriter
	U   *user.User
}
