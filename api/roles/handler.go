// Package roles posts a role and also gets all roles.
package roles

import (
	// "github.com/MerinEREN/iiPackages/api/user"
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/role"
	"github.com/MerinEREN/iiPackages/datastore/roleTypeRole"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// Handler posts a role and returns all the roles from the begining of the kind
// and gets all the roles from the begining of the kind.
func Handler(s *session.Session) {
	// THE CONTROLS BELOVE PREVENT GET REQUEST THAT NECESSARY FOR SELECT SOME SELECT
	// FIELDS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* u, err := user.Get(s)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status == "suspended" {
		log.Printf("Suspended user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are suspended", http.StatusForbidden)
		return
	}
	isAdmin, err := u.IsAdmin(s.Ctx)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	isContentEditor, err := u.IsContentEditor(s.Ctx)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Type != "inHouse" || !(isAdmin || isContentEditor) {
		log.Printf("Unauthorized user %s trying to see "+
			"%s path!!!", u.Email, s.R.URL.Path)
		http.Error(s.W, "You are unauthorized user.", http.StatusUnauthorized)
		return
	} */
	switch s.R.Method {
	case "POST":
		contentID := s.R.FormValue("contentID")
		r := &role.Role{
			ContentID: contentID,
		}
		roleTypesString := s.R.FormValue("types")
		roleTypes := strings.Split(roleTypesString, ",")
		for _, v := range roleTypes {
			if v == "" {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
					"Enpty roleType value or no roleTypes")
				http.Error(s.W, "Enpty roleType value or no roleTypes",
					http.StatusBadRequest)
				return
			}
		}
		kr := datastore.NewKey(s.Ctx, "Role", r.ContentID, 0, nil)
		rtr := &roleTypeRole.RoleTypeRole{
			RoleKey: kr,
		}
		var krtrx []*datastore.Key
		var rtrx roleTypeRole.RoleTypeRoles
		for _, v := range roleTypes {
			krt := datastore.NewKey(s.Ctx, "RoleType", v, 0, nil)
			krtr := datastore.NewIncompleteKey(s.Ctx, "RoleTypeRole", krt)
			krtrx = append(krtrx, krtr)
			rtrx = append(rtrx, rtr)
		}
		opts := new(datastore.TransactionOptions)
		opts.XG = true
		rs := make(role.Roles)
		// "cursor" IS OUT OF USE FOR NOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		// USAGE "s" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!
		err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			rs, err1 = role.PutAndGetMulti(s, r)
			if err1 != nil && err1 != datastore.Done {
				return
			}
			err1 = roleTypeRole.PutMulti(ctx, krtrx, rtrx)
			return
		}, opts)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = rs
		rb.PrevPageURL = "/roles?c=" + crsr.String()
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	default:
		// Handles "GET" requests
		rs, err := role.GetMulti(s.Ctx, nil)
		if err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(rs) == 0 {
			s.W.WriteHeader(http.StatusNoContent)
			return
		}
		rb := new(api.ResponseBody)
		rb.Result = rs
		api.WriteResponseJSON(s, rb)
	}
}
