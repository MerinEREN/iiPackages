/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package role

import (
	"google.golang.org/appengine"
	// "golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"net/http"
)

func Get(r *http.Request, ks []*datastore.Key) (roles Roles, err error) {
	ctx := appengine.NewContext(r)
	err = datastore.GetMulti(ctx, ks, roles)
	return
}

// MAKE THIS MULTI !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Put(r *http.Request) (*datastore.Key, error) {
	ctx := appengine.NewContext(r)
	role := new(Role)
	if r.PostFormValue("roleLang") == "en_us" {
		role.ID = r.PostFormValue("roleValue")
		key := datastore.NewKey(ctx, "Roles", role.ID, 0, nil)
	}
	role.Values[r.PostFormValue("roleLang")] =
		r.PostFormValue(r.PostFormValue(roleValue))
	_, err := datastore.Put(ctx, key, role)
	return err
}
