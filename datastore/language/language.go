/*
Package language "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package language

import (
	"errors"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindLanguage = errors.New("Error while getting language")
)

// GetMulti returns limited entitity from the given cursor.
// If limit is nil default limit will be used.
func GetMulti(s *session.Session, c datastore.Cursor, limit interface{}) (Languages, datastore.Cursor, error) {
	ls := make(Languages)
	// Maybe -LastModied should be the order ctireia if consider UX, think about that.
	q := datastore.NewQuery("Language").Project("Name", "Link").Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	if limit != nil {
		l := limit.(int)
		q = q.Limit(l)
	} else {
		q = q.Limit(20)
	}
	for it := q.Run(s.Ctx); ; {
		l := new(Language)
		k, err := it.Next(l)
		if err == datastore.Done {
			c, err = it.Cursor()
			return ls, c, err
		}
		if err != nil {
			err = ErrFindLanguage
			return nil, c, err
		}
		l.ID = k.StringID()
		ls[l.ID] = l
	}
}

/*
Put "Inside a package, any comment immediately preceding a top-level declaration serves as a
doc comment for that declaration. Every exported (capitalized) name in a program should
have a doc comment.
Doc comments work best as complete sentences, which allow a wide variety of automated
presentations. The first sentence should be a one-sentence summary that starts with the
name being declared."
*/
// Compile parses a regular expression and returns, if successful,
// a Regexp that can be used to match against text.
func Put(s *session.Session, l *Language) (*Language, error) {
	k := datastore.NewKey(s.Ctx, "Language", l.ID, 0, nil)
	var err error
	l.LastModified = time.Now()
	if s.R.Method == "POST" {
		l.Created = time.Now()
		_, err = datastore.Put(s.Ctx, k, l)
	} else {
		// For "PUT" request
		tempL := new(Language)
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			if err1 = datastore.Get(ctx, k, tempL); err1 != nil {
				return
			}
			l.Created = tempL.Created
			_, err1 = datastore.Put(ctx, k, l)
			return
		}, nil)
	}
	return l, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities by the given limit.
func PutAndGetMulti(s *session.Session, c datastore.Cursor, l *Language) (Languages,
	datastore.Cursor, error) {
	ls := make(Languages)
	lNew := new(Language)
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		lNew, err1 = Put(s, l)
		if err1 != nil {
			return
		}
		ls, c, err1 = GetMulti(s, c, 19)
		return
	}, nil)
	ls[lNew.ID] = lNew
	return ls, c, err
}

// GetCount returns language count and an error.
/* func GetCount(s *session.Session) (c int, err error) {
	q := datastore.NewQuery("Language")
	c, err = q.Count(s.Ctx)
	return
} */

// Delete removes the entity by the provided language code and returns an error.
func Delete(s *session.Session, langCode string) error {
	k := datastore.NewKey(s.Ctx, "Language", langCode, 0, nil)
	return datastore.Delete(s.Ctx, k)
}

// DeleteMulti removes the entitys by the provided language code slice
// and returns an error.
func DeleteMulti(s *session.Session, lcx []string) error {
	var kx []*datastore.Key
	for _, v := range lcx {
		kx = append(kx, datastore.NewKey(s.Ctx, "Language", v, 0, nil))
	}
	return datastore.DeleteMulti(s.Ctx, kx)
}
