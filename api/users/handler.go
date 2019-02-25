/*
Package users returns non logged and non deleted users of an account and saves an user
with it's user roles and user tags.
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
)

// Handler returns non logged and non deleted users of an account via account id.
// Also if the request method is POST, puts an user into the User kind
// with account id as parent key and logged user's type as type,
// user roles into the UserRole kind and user tags into the UserTag kind.
func Handler(s *session.Session) {
	var accID string
	if s.R.Header.Get("content-type") == "multipart/form-data" {
		// https://stackoverflow.com/questions/15202448/go-formfile-for-multiple-files
		err := s.R.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile.
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		accID = s.R.Form.Get("accID")
	} else {
		accID = s.R.FormValue("accID")
	}
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
	uLogged, err := userLogged.Get(s)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	var us user.Users
	uNew := new(user.User)
	switch s.R.Method {
	case "POST":
		email := s.R.MultipartForm.Value["email"][0]
		if !valid.IsEmail(email) {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path,
				user.ErrInvalidEmail)
			// NOT SURE ABOUT RETURNED HTTP STATUS !!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Error(s.W, user.ErrInvalidEmail.Error(),
				http.StatusExpectationFailed)
			return
		}
		uNew = &user.User{
			Email: email,
			Type:  uLogged.Type,
		}
		roleIDs := s.R.MultipartForm.Value["roleIDs"]
		tagIDs := s.R.MultipartForm.Value["tagIDs"]
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
				kur := datastore.NewKey(s.Ctx, "UserRole", v, 0, ku)
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
					kut := datastore.NewKey(s.Ctx, "UserTag", v, 0, ku)
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
	default:
		// Handles "GET" requests
	}
	us, err = user.GetProjected(s.Ctx, accKey)
	if err != nil && err != datastore.Done {
		log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		http.Error(s.W, err.Error(), http.StatusInternalServerError)
		return
	}
	// Remove logged and deleted users from the users map.
	for i, v := range us {
		if uLogged.Email == v.Email || v.Status == "deleted" {
			delete(us, i)
		}
	}
	if len(us) == 0 {
		s.W.WriteHeader(http.StatusNoContent)
		return
	}
	if s.R.Method == "POST" {
		us[uNew.ID] = uNew
	}
	rb := new(api.ResponseBody)
	rb.Result = us
	if s.R.Method == "POST" {
		s.W.Header().Set("Content-Type", "application/json")
		s.W.WriteHeader(http.StatusCreated)
		api.WriteResponse(s, rb)
	} else {
		api.WriteResponseJSON(s, rb)
	}
}
