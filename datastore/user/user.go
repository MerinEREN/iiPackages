/*
Package user has user struct type and datastore query functions.
*/
package user

import (
	"errors"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/userRole"
	"github.com/MerinEREN/iiPackages/session"
	valid "github.com/asaskevich/govalidator"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrEmailExist      = errors.New("email exist")
	ErrEmailNotExist   = errors.New("email not exist")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrPutUser         = errors.New("error while putting user into the datastore")
)

// IsAdmin gets user's role keys first. Than returns encoded content key from role key's
// stringID and converts it to content key. And finaly, gets content's stringID from it
// and compares.
// Than returns a boolean and an error.
func (u *User) IsAdmin(ctx context.Context) (bool, error) {
	k, err := datastore.DecodeKey(u.ID)
	if err != nil {
		return false, err
	}
	krx, err := userRole.GetKeysUserOrRole(ctx, k)
	if err != nil && err != datastore.Done {
		return false, err
	}
	var encodedContentKey string
	kc := new(datastore.Key)
	for _, v := range krx {
		encodedContentKey = v.StringID()
		kc, err = datastore.DecodeKey(encodedContentKey)
		if err != nil {
			return false, err
		}
		if kc.StringID() == "admin" {
			return true, nil
		}
	}
	return false, nil
}

// IsContentEditor gets user's role keys first. Than returns encoded content key
// from role key's stringID and converts it to content key.
// And finaly, gets content's stringID from it and compares.
// Than returns a boolean and an error.
func (u *User) IsContentEditor(ctx context.Context) (bool, error) {
	k, err := datastore.DecodeKey(u.ID)
	if err != nil {
		return false, err
	}
	krx, err := userRole.GetKeysUserOrRole(ctx, k)
	if err != nil && err != datastore.Done {
		return false, err
	}
	var encodedContentKey string
	kc := new(datastore.Key)
	for _, v := range krx {
		encodedContentKey = v.StringID()
		kc, err = datastore.DecodeKey(encodedContentKey)
		if err != nil {
			return false, err
		}
		if kc.StringID() == "contentEditor" {
			return true, nil
		}
	}
	return false, nil
}

// CreateWithAccount creates an account and an user entities for the new user
// in a transaction.
// And returns account, user, user key and an error.
func CreateWithAccount(s *session.Session) (acc *account.Account, u *User,
	k *datastore.Key, err error) {
	email := s.U.Email
	// Email validation control not necessary for now.
	if !valid.IsEmail(email) {
		err = ErrInvalidEmail
		return nil, nil, nil, err
	}
	// CAHANGE THIS CONTROL AND ALLOW SPECIAL CHARACTERS !!!!!!!!!!!!!!!!!!!!!!
	/* if !valid.IsAlphanumeric(password) {
		err = usr.InvalidPassword
		return nil, nil, err
	} */
	// UUID has no functionallity for now.
	u4 := new(uuid.UUID)
	if u4, err = uuid.NewV4(); err != nil {
		return nil, nil, nil, err
	}
	UUID := u4.String()
	acc = &account.Account{
		Registered:   time.Now(),
		LastModified: time.Now(),
	}
	ka := datastore.NewKey(s.Ctx, "Account", UUID, 0, nil)
	// USAGE "s" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!!!!!!!!!
	err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if ka, err1 = datastore.Put(ctx, ka, acc); err1 != nil {
			return
		}
		// DELETE THIS WHEN U STORE ONLY KEYS INTO MEMCACHE !!!!!!!!!!!!!!!!!!!!!!!
		acc.ID = ka.Encode()
		u, k, err1 = New(s, email, ka)
		// DELETE THIS WHEN U STORE ONLY KEYS INTO MEMCACHE !!!!!!!!!!!!!!!!!!!!!!!
		u.ID = k.Encode()
		return
	}, nil)
	return
}

// New creates a new user and returns it and its key if not exist, or returns an error.
// Also sets new user's default role as admin to "userRole" kind in a transaction.
func New(s *session.Session, email string, pk *datastore.Key) (
	u *User, k *datastore.Key, err error) {
	var c int
	c, err = datastore.NewQuery("User").Filter("Email =", email).Count(s.Ctx)
	if err != nil {
		return
	} else if c > 0 {
		err = ErrEmailExist
	} else {
		u = &User{
			Email:        email,
			IsActive:     true,
			Created:      time.Now(),
			LastModified: time.Now(),
			// Password:     GetHmac(password),
		}
		k = datastore.NewKey(s.Ctx, "User", email, 0, pk)
		kc := datastore.NewKey(s.Ctx, "Content", "admin", 0, nil)
		kr := datastore.NewKey(s.Ctx, "Role", kc.Encode(), 0, nil)
		ur := &userRole.UserRole{
			RoleKey: kr,
		}
		kur := datastore.NewIncompleteKey(s.Ctx, "UserRole", k)
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			if k, err1 = datastore.Put(ctx, k, u); err1 != nil {
				return
			}
			err1 = userRole.Put(ctx, kur, ur)
			return
		}, nil)
		u.ID = k.Encode()
	}
	return
}

// Get returns the entity and an error via entity's key.
func Get(ctx context.Context, k *datastore.Key) (*User, error) {
	u := new(User)
	err := datastore.Get(ctx, k, u)
	if err != nil {
		return nil, err
	}
	u.ID = k.Encode()
	return u, err
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

// GetProjected returns limited entities from the given cursor
// with thumbnail enough properties via account key, the updated cursor and an error.
// If limit is nil default limit will be used.
func GetProjected(ctx context.Context, pk *datastore.Key, crsr datastore.Cursor,
	limit interface{}) (Users, datastore.Cursor, error) {
	us := make(Users)
	q := datastore.NewQuery("User").
		Ancestor(pk).
		Order("-Created").
		Project("Name.First", "Name.Last", "Email", "Link")
	if crsr.String() != "" {
		q = q.Start(crsr)
	}
	if limit != nil {
		l := limit.(int)
		q = q.Limit(l)
	} else {
		q = q.Limit(10)
	}
	for it := q.Run(ctx); ; {
		u := new(User)
		k, err := it.Next(u)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return us, crsr, err
		}
		if err != nil {
			return nil, crsr, err
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

// Put adds to or modifies an entity in the kind according to request method.
func Put(s *session.Session, u *User, pk *datastore.Key) (*datastore.Key, *User, error) {
	k := new(datastore.Key)
	var err error
	if s.R.Method == "POST" {
		var c int
		c, err = datastore.NewQuery("User").Filter("Email =", u.Email).Count(s.Ctx)
		if err != nil {
			return nil, nil, err
		} else if c > 0 {
			err = ErrEmailExist
			return nil, nil, err
		} else {
			k = datastore.NewKey(s.Ctx, "User", u.Email, 0, pk)
			u.ID = k.Encode()
			u.Created = time.Now()
		}
	} else if s.R.Method == "PUT" {
		k, err = datastore.DecodeKey(u.ID)
		if err != nil {
			return nil, nil, err
		}
		tempU := new(User)
		if err = datastore.Get(s.Ctx, k, tempU); err != nil {
			return nil, nil, err
		}
		u.Created = tempU.Created
	}
	u.LastModified = time.Now()
	k, err = datastore.Put(s.Ctx, k, u)
	return k, u, err
}
