/*
Package user "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package user

import (
	// "fmt"
	// "google.golang.org/appengine"
	// "google.golang.org/appengine/user"
	// "github.com/MerinEREN/iiPackages/cookie"
	// "github.com/MerinEREN/iiPackages/datastore/role"
	// valid "github.com/asaskevich/govalidator"
	// "github.com/MerinEREN/iiPackages/datastore/tag"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	// "io"
	// "log"
	// "net/http"
	"errors"
	"time"
)

// Errors
var (
	ErrEmailNotExist   = errors.New("email not exist")
	ErrEmailExist      = errors.New("email exist")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrPutUser         = errors.New("error while putting user into the datastore")
	ErrFindUser        = errors.New("error while checking email existincy")
)

// IsAdmin "Exported functions should have a comment"
func (u *User) IsAdmin() bool {
	for _, r := range u.Roles {
		if r == "admin" {
			return true
		}
	}
	return false
}

// IsContentEditor "Exported functions should have a comment"
func (u *User) IsContentEditor() bool {
	for _, r := range u.Roles {
		if r == "contentEditor" {
			return true
		}
	}
	return false
}

// New "Exported functions should have a comment"
func New(ctx context.Context, parentKey *datastore.Key, email, role string) (u *User,
	key *datastore.Key, err error) {
	var roles []string
	roles = append(roles, role)
	u, key, err = Get(ctx, email, parentKey)
	switch err {
	case nil:
		err = ErrEmailExist
	case datastore.ErrNoSuchEntity:
		u = &User{
			ID:           email,
			Email:        email,
			Roles:        roles,
			IsActive:     true,
			Registered:   time.Now(),
			LastModified: time.Now(),
			// Password:     GetHmac(password),
		}
		_, err = datastore.Put(ctx, key, u)
	}
	return
}

// Get "Exported functions should have a comment"
func Get(ctx context.Context, email string, parentKey *datastore.Key) (*User,
	*datastore.Key, error) {
	// BUG !!!!! If i made this function as naked return "u.ID = email" fails
	// because of "u" !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	u := new(User)
	k := new(datastore.Key)
	var err error
	k = datastore.NewKey(ctx, "User", email, 0, parentKey)
	err = datastore.Get(ctx, k, u)
	u.ID = email
	return u, k, err
}

// GetWithEmail "Exported functions should have a comment"
func GetWithEmail(ctx context.Context, email string) (*User, *datastore.Key, error) {
	u := new(User)
	q := datastore.NewQuery("User").Filter("Email =", email)
	it := q.Run(ctx)
	// BUG !!!!! If i made this function as naked return "it.Next" fails because of "u"
	k, err := it.Next(u)
	if err == datastore.Done {
		return nil, nil, err
	}
	if err != nil {
		err = ErrFindUser
		return nil, nil, err
	}
	u.ID = k.StringID()
	return u, k, nil
}

// GetKey "Exported functions should have a comment"
func GetKey(ctx context.Context, email string) (k *datastore.Key, err error) {
	q := datastore.NewQuery("User").Filter("Email =", email).KeysOnly()
	it := q.Run(ctx)
	k, err = it.Next(nil)
	if err == datastore.Done {
		return
	}
	if err != nil {
		err = ErrFindUser
		return
	}
	return
}

// GetAllKeys "Exported functions should have a comment"
func GetAllKeys(ctx context.Context, tagIDs []*datastore.Key) (ks []*datastore.Key, err error) {
	// Check what returns when done while using GetAll() !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	ks, err = datastore.NewQuery("User").
		Filter("Tags =", tagIDs).
		KeysOnly().
		GetAll(ctx, nil)
	return
}

// Exist "Exported functions should have a comment"
func Exist(ctx context.Context, email string) (c int, err error) {
	c, err = datastore.NewQuery("User").Filter("Email =", email).Count(ctx)
	return
}

// GetTagIDs "Exported functions should have a comment"
func GetTagIDs(ctx context.Context, email string) ([]*datastore.Key, error) {
	u := new(User)
	q := datastore.NewQuery("User").Filter("Email =", email).Project("TagIDs")
	it := q.Run(ctx)
	_, err := it.Next(u)
	if err == datastore.Done {
		return nil, err
	}
	if err != nil {
		err = ErrFindUser
		return nil, err
	}
	return u.TagIDs, nil
}
