/*
Package session contains session struckt and it's methods.
*/
package session

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
	"net/http"
)

// Session struct has Context, Request, ResponseWriter and logged User fields.
type Session struct {
	Ctx context.Context
	R   *http.Request
	W   http.ResponseWriter
	U   *user.User
}
