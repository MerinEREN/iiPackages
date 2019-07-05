/*
Package user has user struct type and datastore query functions.
*/
package user

import (
	"errors"
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/roleUser"
	"github.com/MerinEREN/iiPackages/datastore/tagUser"
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
	krx, err := roleUser.GetKeysByUserOrRoleKey(ctx, k)
	if err != nil && err != datastore.Done {
		return false, err
	}
	var encodedContextKey string
	kc := new(datastore.Key)
	for _, v := range krx {
		encodedContextKey = v.StringID()
		kc, err = datastore.DecodeKey(encodedContextKey)
		if err != nil {
			return false, err
		}
		if kc.StringID() == "admin" {
			return true, nil
		}
	}
	return false, nil
}

// IsContextEditor gets user's role keys first. Than returns encoded content key
// from role key's stringID and converts it to content key.
// And finaly, gets content's stringID from it and compares.
// Than returns a boolean and an error.
func (u *User) IsContextEditor(ctx context.Context) (bool, error) {
	k, err := datastore.DecodeKey(u.ID)
	if err != nil {
		return false, err
	}
	krx, err := roleUser.GetKeysByUserOrRoleKey(ctx, k)
	if err != nil && err != datastore.Done {
		return false, err
	}
	var encodedContextKey string
	kc := new(datastore.Key)
	for _, v := range krx {
		encodedContextKey = v.StringID()
		kc, err = datastore.DecodeKey(encodedContextKey)
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
// Also sets new user's default role as admin to "roleUser" kind in a transaction.
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
			Created:      time.Now(),
			LastModified: time.Now(),
			// Password:     GetHmac(password),
		}
		k = datastore.NewKey(s.Ctx, "User", email, 0, pk)
		kc := datastore.NewKey(s.Ctx, "Context", "admin", 0, nil)
		kr := datastore.NewKey(s.Ctx, "Role", kc.Encode(), 0, nil)
		ru := &roleUser.RoleUser{
			RoleKey: kr,
		}
		kru := datastore.NewKey(s.Ctx, "RoleUser", kr.Encode(), 0, k)
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			if k, err1 = datastore.Put(ctx, k, u); err1 != nil {
				return
			}
			_, err1 = datastore.Put(ctx, kru, ru)
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

// GetKeysByParent returns an error and users keys as a slice via account key.
func GetKeysByParent(ctx context.Context, pk *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("User")
	q = q.Ancestor(pk).
		KeysOnly()
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

/*
GetKeysByParentOrdered returns users keys as a slice via account key
as descended order of the "Created" property and an error.
*/
func GetKeysByParentOrdered(ctx context.Context, pk *datastore.Key) ([]*datastore.Key, error) {
	q := datastore.NewQuery("User")
	q = q.
		Ancestor(pk).
		Order("-Created").
		KeysOnly()
	return q.GetAll(ctx, nil)
}

// GetProjected returns thumbnail enough properties via account key and an error.
func GetProjected(ctx context.Context, pk *datastore.Key) (Users, error) {
	us := make(Users)
	q := datastore.NewQuery("User").
		Ancestor(pk).
		Order("-Created").
		Project("Name.First", "Name.Last", "Email", "Link", "Status")
	for it := q.Run(ctx); ; {
		u := new(User)
		k, err := it.Next(u)
		if err == datastore.Done {
			return us, err
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

// Put adds to or modifies an entity in the kind according to existance.
func Put(s *session.Session, u *User, pk *datastore.Key) (*datastore.Key, *User, error) {
	k := new(datastore.Key)
	c, err := datastore.NewQuery("User").Filter("Email =", u.Email).Count(s.Ctx)
	if err != nil {
		return nil, nil, err
	}
	if c == 0 {
		// Adding an entity first time.
		k = datastore.NewKey(s.Ctx, "User", u.Email, 0, pk)
		u.ID = k.Encode()
		u.Created = time.Now()
	} else {
		if s.R.Method == "POST" {
			// Adding a removed entity again.
			k = datastore.NewKey(s.Ctx, "User", u.Email, 0, pk)
			u.ID = k.Encode()
		} else {
			// Updating an egzisting entity.
			k, err = datastore.DecodeKey(u.ID)
			if err != nil {
				return nil, nil, err
			}
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

// UpdateStatus set's user status to given value "v" by given encoded user key "ek"
// and returns an error.
func UpdateStatus(ctx context.Context, ek, v string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	u := new(User)
	if err = datastore.Get(ctx, k, u); err != nil {
		return err
	}
	u.Status = v
	_, err = datastore.Put(ctx, k, u)
	return err
}

// Delete sets user's status to "deleted"
// and removes user's roles and tags in a transaction.
// Also, returns an error.
func Delete(ctx context.Context, ek string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	krux, err := roleUser.GetKeys(ctx, k)
	if err != datastore.Done {
		return err
	}
	ktux, err := tagUser.GetKeys(ctx, k)
	if err != datastore.Done {
		return err
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		if err1 = datastore.DeleteMulti(ctx, krux); err1 != nil {
			return
		}
		if err1 = datastore.DeleteMulti(ctx, ktux); err1 != nil {
			return
		}
		err1 = UpdateStatus(ctx, ek, "deleted")
		return
	}, opts)
	return err
}
