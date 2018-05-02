/*
Package account "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package account

import (
	// "fmt"
	"errors"
	"golang.org/x/net/context"
	// "google.golang.org/appengine"
	// "cookie"
	usr "github.com/MerinEREN/iiPackages/datastore/user"
	valid "github.com/asaskevich/govalidator"
	"google.golang.org/appengine/user"
	// "github.com/nu7hatch/gouuid"
	"github.com/nu7hatch/gouuid"
	"google.golang.org/appengine/datastore"
	// "log"
	// "net/http"
	"time"
)

// Errors
var (
	ErrPutAccount  = errors.New("error while putting account into the datastore")
	ErrFindAccount = errors.New("error while getting account")
)

// CreateAccountAndUser creates an account for the new user in a transaction.
/*
Inside a package, any comment immediately preceding a top-level declaration serves as a
doc comment for that declaration. Every exported (capitalized) name in a program should
have a doc comment.
Doc comments work best as complete sentences, which allow a wide variety of automated
presentations. The first sentence should be a one-sentence summary that starts with the
name being declared.
*/
// Compile parses a regular expression and returns, if successful,
// a Regexp that can be used to match against text.
func CreateAccountAndUser(ctx context.Context) (acc *Account, u *usr.User,
	uK *datastore.Key, err error) {
	ug := user.Current(ctx)
	email := ug.Email
	// Email validation control not necessary actually.
	if !valid.IsEmail(email) {
		err = usr.ErrInvalidEmail
		return nil, nil, nil, err
	}
	// CAHANGE THIS CONTROL AND ALLOW SPECIAL CHARACTERS !!!!!!!!!!!!!!!!!!!!!!
	/* if !valid.IsAlphanumeric(password) {
		err = usr.InvalidPassword
		return nil, nil, err
	} */
	u4 := new(uuid.UUID)
	if u4, err = uuid.NewV4(); err != nil {
		return nil, nil, nil, err
	}
	UUID := u4.String()
	acc = &Account{
		ID:           UUID,
		Registered:   time.Now(),
		LastModified: time.Now(),
	}
	key := datastore.NewKey(ctx, "Account", UUID, 0, nil)
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		if _, err1 = datastore.Put(ctx, key, acc); err1 != nil {
			return
		}
		u, uK, err1 = usr.New(ctx, key, email, "admin")
		return
	}, nil)
	return
}

// Get "Exported functions should have a comment"
func Get(ctx context.Context, k interface{}) (*Account, error) {
	acc := new(Account)
	var err error
	switch v := k.(type) {
	case string:
		// Do some projection here if needed
		q := datastore.NewQuery("Account").
			Filter("KeyName =", v)
		it := q.Run(ctx)
		_, err = it.Next(acc)
		acc.ID = v
	case *datastore.Key:
		err = datastore.Get(ctx, v, acc)
		acc.ID = v.StringID()
	}
	if err == datastore.Done {
		return acc, err
	}
	if err != nil {
		err = ErrFindAccount
		return nil, err
	}
	return acc, nil
}

// Delete "Exported functions should have a comment"
func Delete(ctx context.Context, k *datastore.Key) error {
	err := datastore.Delete(ctx, k)
	return err
}

/* func AddTags(s ...string) bool {
	return
} */
