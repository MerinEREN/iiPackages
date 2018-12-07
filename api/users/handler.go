/*
Package users returns only non adimin users of an account and saves user user roles
and user tags.
*/
package users

import (
	"github.com/MerinEREN/iiPackages/api"
	userLogged "github.com/MerinEREN/iiPackages/api/user"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/datastore/userRole"
	"github.com/MerinEREN/iiPackages/datastore/userTag"
	"github.com/MerinEREN/iiPackages/session"
	valid "github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strings"
)

// Handler returns the only non logged users of an account via account id
// with limited properties to list in User Settings page.
// Also puts a user into the User kind with account id as parent key
// and logged user's type as type,
// user roles into the UserRole kind and user tags into the UserTag kind
// if the request method is POST.
func Handler(s *session.Session) {
	accID := s.R.FormValue("accID")
	// CHECK THIS KONTROL BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	if accID == "" {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, "No account ID")
		http.Error(s.W, "No account ID", http.StatusBadRequest)
		return
	}
	accKey, err := datastore.DecodeKey(accID)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	var crsr datastore.Cursor
	var us user.Users
	var count interface{}
	uLogged := new(user.User)
	uNew := new(user.User)
	switch s.R.Method {
	case "POST":
		email := s.R.FormValue("email")
		if !valid.IsEmail(email) {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
				user.ErrInvalidEmail)
			// NOT SURE ABOUT RETURNED HTTP STATUS !!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(s.W, user.ErrInvalidEmail.Error(),
				http.StatusExpectationFailed)
			return
		}
		uLogged, err = userLogged.Get(s)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		uNew = &user.User{
			Email:    email,
			IsActive: true,
			Type:     uLogged.Type,
		}
		roleIDsString := s.R.FormValue("roleIDs")
		tagIDsString := s.R.FormValue("tagIDs")
		roleIDs := strings.Split(roleIDsString, ",")
		for _, v := range roleIDs {
			if v == "" {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
					"Enpty roleID value or no roleIDs")
				http.Error(s.W, "Enpty roleID value or no roleIDs",
					http.StatusBadRequest)
				return
			}
		}
		tagIDs := strings.Split(tagIDsString, ",")
		ku := new(datastore.Key)
		var urx userRole.UserRoles
		var utx userTag.UserTags
		var kurx []*datastore.Key
		var kutx []*datastore.Key
		opts := new(datastore.TransactionOptions)
		opts.XG = true
		// USAGE "s" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			ku, uNew, err1 = user.Put(s, uNew, accKey)
			if err1 != nil {
				return
			}
			for _, v := range roleIDs {
				ur := new(userRole.UserRole)
				kr, err := datastore.DecodeKey(v)
				if err != nil {
					return
				}
				ur.RoleKey = kr
				urx = append(urx, ur)
				kur := datastore.NewIncompleteKey(s.Ctx, "UserRole", ku)
				kurx = append(kurx, kur)
			}
			if len(kurx) > 0 {
				_, err1 = datastore.PutMulti(ctx, kurx, urx)
				if err1 != nil {
					return
				}
			}
			for _, v := range tagIDs {
				if v != "" {
					ut := new(userTag.UserTag)
					kt, err := datastore.DecodeKey(v)
					if err != nil {
						return
					}
					ut.TagKey = kt
					utx = append(utx, ut)
					kut := datastore.NewIncompleteKey(s.Ctx, "UserTag", ku)
					kutx = append(kutx, kut)
				}
			}
			if len(kutx) > 0 {
				_, err1 = datastore.PutMulti(ctx, kutx, utx)
				if err1 != nil {
					return
				}
			}
			return
		}, opts)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		// Default limit +1 because of the admin user removal.
		// Not 11 but the default value 10
		// beacause of the newUser addition to the response data.
		count = nil
	default:
		// Handles "GET" requests
		crsr, err = datastore.DecodeCursor(s.R.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		// if first pagination request.
		if crsr.String() == "" {
			// Default limit +1 because of the admin user removal.
			count = 11
		} else {
			count = nil
		}
	}
	// Get the entities from the begining by reseted cursor if the method is Post
	// or get from the given cursor.
	us, crsr, err = user.GetProjected(s.Ctx, accKey, crsr, count)
	if err != nil && err != datastore.Done {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	// Remove admin user from the users map.
	/* var isAdmin bool
	for i, v := range us {
		isAdmin, err = v.IsAdmin(s.Ctx)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		if isAdmin {
			delete(us, i)
			break
		}
	} */
	// Remove logged user from the users map.
	for i, v := range us {
		if uLogged.Email == v.Email {
			delete(us, i)
			break
		}
	}
	if len(us) == 0 {
		s.W.WriteHeader(http.StatusNoContent)
		return
	}
	if s.R.Method == "POST" {
		us[uNew.ID] = uNew
	}
	// Returns only non admin users of the account.
	rb := new(api.ResponseBody)
	rb.Result = us
	rb.PrevPageURL = "/users?c=" + crsr.String() + "&accID=" + accID
	if s.R.Method == "POST" {
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	} else {
		api.WriteResponseJSON(s, rb)
	}
}
