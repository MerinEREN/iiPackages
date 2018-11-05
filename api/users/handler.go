// Package users returns only non adimin users of an account and saves user and user tags.
package users

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/datastore/userTag"
	"github.com/MerinEREN/iiPackages/session"
	valid "github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
)

// Handler returns the only non adimin users of an account via account id
// with limited properties to list in User Settings page.
// Also puts a user into the User kind with account id as parent key
// and user tags into the UserTag kind if the request method is POST.
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
	uNew := new(user.User)
	var crsr datastore.Cursor
	var us user.Users
	var count interface{}
	switch s.R.Method {
	case "POST":
		// GET ROLES TOO !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		email := s.R.FormValue("email")
		if !valid.IsEmail(email) {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
				user.ErrInvalidEmail)
			// NOT SURE ABOUT RETURNED HTTP STATUS !!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(s.W, user.ErrInvalidEmail.Error(),
				http.StatusExpectationFailed)
			return
		}
		tagIDs := s.R.FormValue("tagIDs")
		uk := new(datastore.Key)
		utx := make(userTag.UserTags, len(tagIDs))
		utkx := make([]*datastore.Key, len(tagIDs))
		opts := new(datastore.TransactionOptions)
		opts.XG = true
		err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			uNew, uk, err1 = user.New(s, email, nil, accKey)
			if err1 != nil {
				return
			}
			for _, v := range tagIDs {
				ut := new(userTag.UserTag)
				// tk, err := datastore.DecodeKey(v)
				tk, err := datastore.DecodeKey(string(v))
				if err != nil {
					return
				}
				ut.TagKey = tk
				utx = append(utx, ut)
				utk := datastore.NewKey(s.Ctx, "UserTag", string(v), 0, uk)
				utkx = append(utkx, utk)
			}
			if len(utkx) > 0 {
				_, err1 = datastore.PutMulti(ctx, utkx, utx)
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
	// Or get from the given cursor.
	us, crsr, err = user.GetProjected(s.Ctx, accKey, crsr, count)
	if err != nil && err != datastore.Done {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	// Remove admin user from the users map
	brk := false
	for i, v := range us {
		if brk {
			break
		}
		for _, v2 := range v.Roles {
			if v2 == "admin" {
				delete(us, i)
				brk = true
				break
			}
		}
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
