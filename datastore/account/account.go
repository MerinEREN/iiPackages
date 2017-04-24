/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
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
	"google.golang.org/appengine/datastore"
	// "log"
	// "net/http"
	"strconv"
	"strings"
	"time"
)

// Errors
var (
	ErrPutAccount  = errors.New("Error while putting account into the datastore.")
	ErrFindAccount = errors.New("Error while getting account.")
)

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
func Create(ctx context.Context) (acc *Account, u *usr.User, uK *datastore.Key,
	err error) {
	// CAHANGE THIS CONTROL AND ALLOW SPECIAL CHARACTERS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* if !valid.IsAlphanumeric(password) {
		err = usr.InvalidPassword
		return
	} */
	var accKeyName string
	key := new(datastore.Key)
	q := datastore.NewQuery("Account").
		Order("-Registered").
		KeysOnly()
	it := q.Run(ctx)
	key, err = it.Next(nil)
	if err != nil {
		if err == datastore.Done {
			accKeyName = "Account_1"
		} else {
			return
		}
	} else {
		s := strings.SplitAfter(key.StringID(), "_")
		var i int
		i, err = strconv.Atoi(s[1])
		if err != nil {
			return
		}
		i = i + 1
		accKeyName = s[0] + strconv.Itoa(i)
	}
	acc = &Account{
		ID:           accKeyName,
		Registered:   time.Now(),
		LastModified: time.Now(),
	}
	key = datastore.NewKey(ctx, "Account", accKeyName, 0, nil)
	key, err = datastore.Put(ctx, key, acc)
	if err != nil {
		return
	}
	ug := user.Current(ctx)
	email := ug.Email
	// Email validation control not necessary actually.
	if !valid.IsEmail(email) {
		err = usr.ErrInvalidEmail
		return
	}
	u, uK, err = usr.New(ctx, key, email, "admin")
	if err != nil {
		if errD := Delete(ctx, key); errD != nil {
			err = errors.New(err.Error() + errD.Error())
		}
		return
	}
	return
}

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

func Delete(ctx context.Context, k *datastore.Key) error {
	err := datastore.Delete(ctx, k)
	return err
}

/* func AddTags(s ...string) bool {
	return
} */
