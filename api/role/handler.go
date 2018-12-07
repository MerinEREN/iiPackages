// Package role deletes a role by the given id as encoded datastore key.
package role

import (
	"github.com/MerinEREN/iiPackages/datastore/role"
	"github.com/MerinEREN/iiPackages/datastore/roleTypeRole"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// Handler handles delete request to delete a role by the given id as encoded datastore key.
func Handler(s *session.Session) {
	ID := strings.Split(s.R.URL.Path, "/")[2]
	switch s.R.Method {
	case "DELETE":
		if ID == "" {
			log.Printf("Path: %s, Error: no role ID to delete\n", s.R.URL.Path)
			http.Error(s.W, "No role ID to delete", http.StatusBadRequest)
			return
		}
		k, err := datastore.DecodeKey(ID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		opts := new(datastore.TransactionOptions)
		opts.XG = true
		krtrx, err := roleTypeRole.GetKeys(s.Ctx, k)
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			err1 = roleTypeRole.DeleteMulti(ctx, krtrx)
			if err1 != nil {
				return
			}
			err = role.Delete(ctx, k)
			return
		}, opts)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		// Handles "PUT" requests
	}
}
