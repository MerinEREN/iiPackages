/*
Package user has user struct type and datastore query functions.
*/
package user

import (
	"errors"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/session"
	valid "github.com/asaskevich/govalidator"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	// ErrEmailExist      = errors.New("email exist")
	ErrEmailNotExist   = errors.New("email not exist")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrPutUser         = errors.New("error while putting user into the datastore")
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

// CreateWithAccount creates an account and an user entities for the new user
// in a transaction.
// And returns account, user, user key and an error.
func CreateWithAccount(s *session.Session) (acc *account.Account, u *User,
	uKey *datastore.Key, err error) {
	// Email validation control not necessary actually.
	if !valid.IsEmail(s.U.Email) {
		err = ErrInvalidEmail
		return nil, nil, nil, err
	}
	// CAHANGE THIS CONTROL AND ALLOW SPECIAL CHARACTERS !!!!!!!!!!!!!!!!!!!!!!
	/* if !valid.IsAlphanumeric(password) {
		err = usr.InvalidPassword
		return nil, nil, err
	} */
	// UUID is not necessary actually.
	u4 := new(uuid.UUID)
	if u4, err = uuid.NewV4(); err != nil {
		return nil, nil, nil, err
	}
	UUID := u4.String()
	acc = &account.Account{
		ID:           UUID,
		Registered:   time.Now(),
		LastModified: time.Now(),
	}
	k := datastore.NewKey(s.Ctx, "Account", UUID, 0, nil)
	err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if k, err1 = datastore.Put(ctx, k, acc); err1 != nil {
			return
		}
		u, uKey, err1 = New(s, k)
		return
	}, nil)
	return
}

// New creates a new user and returns it and its key if not exist or returns an error.
func New(s *session.Session, pk *datastore.Key) (u *User, k *datastore.Key, err error) {
	var roles []string
	email := s.U.Email
	roles = append(roles, "admin")
	// Existance check is not necessary.
	/* var c int
	c, err = datastore.NewQuery("User").Filter("Email =", email).Count(s.Ctx)
	if err != nil {
		return
	} else if c > 0 {
		err = ErrEmailExist
	} else { */
	u = &User{
		Email:        email,
		Roles:        roles,
		IsActive:     true,
		Registered:   time.Now(),
		LastModified: time.Now(),
		// Password:     GetHmac(password),
	}
	k = datastore.NewKey(s.Ctx, "User", email, 0, pk)
	k, err = datastore.Put(s.Ctx, k, u)
	u.ID = k.Encode()
	// }
	return
}

// Get returns the entity and an error via entity's key.
func Get(s *session.Session, k *datastore.Key) (u *User, err error) {
	err = datastore.Get(s.Ctx, k, u)
	u.ID = k.Encode()
	return
}

// GetUsersKeysViaParent returns an error and users keys as a slice via account key.
func GetUsersKeysViaParent(ctx context.Context, pk *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("User").Ancestor(pk).KeysOnly()
	for it := q.Run(ctx); ; {
		k, err := it.Next(nil)
		if err == datastore.Done {
			return kx, err
		}
		if err != nil {
			return nil, err
		}
		kx = append(kx, k)
	}
}

// GetProjected return users with thumbnail enough properties via account key
// and an error.
func GetProjected(ctx context.Context, pk *datastore.Key) (Users, error) {
	us := make(Users)
	q := datastore.NewQuery("User")
	q = q.Ancestor(pk)
	q = q.Project("Name", "Link")
	for it := q.Run(ctx); ; {
		u := new(User)
		k, err := it.Next(u)
		if err == datastore.Done {
			return us, nil
		}
		if err != nil {
			return nil, err
		}
		u.ID = k.Encode()
		us[u.ID] = u
	}
}

// GetViaEmailAndParent returns user, user key and an error via email and parent key.
/* func GetViaEmailAndParent(ctx context.Context, email string, pk *datastore.Key) (*User,
	*datastore.Key, error) {
	// BUG !!!!! If i made this function as naked return "u.ID = email" fails
	// because of "u" !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	u := new(User)
	k := datastore.NewKey(ctx, "User", email, 0, pk)
	err := datastore.Get(ctx, k, u)
	return u, k, err
} */

// GetViaEmail returns user and an error via user email only.
func GetViaEmail(s *session.Session) (*User, *datastore.Key, error) {
	u := new(User)
	q := datastore.NewQuery("User").Filter("Email =", s.U.Email)
	it := q.Run(s.Ctx)
	// BUG !!!!! If i made this function as naked return "it.Next" fails because of "u"
	k, err := it.Next(u)
	if err == nil {
		u.ID = k.Encode()
	}
	return u, k, err
}

// GetKeyViaEmail returns user key via email only.
func GetKeyViaEmail(s *session.Session) (k *datastore.Key, err error) {
	q := datastore.NewQuery("User").Filter("Email =", s.U.Email).KeysOnly()
	it := q.Run(s.Ctx)
	k, err = it.Next(nil)
	return
}
