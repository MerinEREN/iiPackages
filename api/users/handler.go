/*
Package users returns users with limited properties to list in User Settings page.
*/
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

// Handler returns users via account id.
// Also puts a user into the User kind and user tags into the UserTag kind
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
	rb := new(api.ResponseBody)
	var us user.Users
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
		/* bs, err := ioutil.ReadAll(s.R.Body)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		p := new(page.Page)
		err = json.Unmarshal(bs, p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		} */
		// Using 'decoder' is an alternative and can be used if response body has
		// more than one json object.
		// Otherwise don't use it, because it has performance disadvantages
		// compared to first solution.
		/*decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(p)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		} */
		utx := make(userTag.UserTags, len(tagIDs))
		utkx := make([]*datastore.Key, len(tagIDs))
		for _, v := range tagIDs {
			ut := new(userTag.UserTag)
			// tk, err := datastore.DecodeKey(v)
			tk, err := datastore.DecodeKey(string(v))
			if err != nil {
				return
			}
			ut.TagKey = tk
			utx = append(utx, ut)
			utk := datastore.NewIncompleteKey(s.Ctx, "UserTag", nil)
			utkx = append(utkx, utk)
		}
		uNew := new(user.User)
		uk := new(datastore.Key)
		opts := new(datastore.TransactionOptions)
		opts.XG = true
		err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			uNew, uk, err1 = user.New(s, email, nil, accKey)
			if err1 != nil {
				return
			}
			for i := range utkx {
				utx[i].UserKey = uk
			}
			_, err1 = datastore.PutMulti(ctx, utkx, utx)
			if err1 != nil {
				return
			}
			return
		}, opts)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		// Reset the cursor and get the entities from the begining.
		var crsr datastore.Cursor
		us, crsr, err = user.GetProjected(s.Ctx, accKey, crsr, 9)
		if err != nil && err != datastore.Done {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		us[uNew.ID] = uNew
		rb.Reset = true
		rb.PrevPageURL = "/users?c=" + crsr.String() + "&accID=" + accID
		s.W.WriteHeader(http.StatusCreated)
	default:
		// Hendles "GET" requests
		// And returns only non admin users of the account
		crsr, err := datastore.DecodeCursor(s.R.FormValue("c"))
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		us, crsr, err = user.GetProjected(s.Ctx, accKey, crsr, nil)
		if err != nil {
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
		rb.PrevPageURL = "/users?c=" + crsr.String() + "&accID=" + accID
	}
	rb.Result = us
	api.WriteResponse(s, rb)
}
