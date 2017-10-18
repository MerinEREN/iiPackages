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
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindLanguage = errors.New("Error while getting language")
)

// GetMulti "Exported functions should have a comment"
func GetMulti(s *session.Session, c datastore.Cursor) (Languages, datastore.Cursor, error) {
	ls := make(Languages)
	// Maybe -LastModied should be the order ctireia if consider UX, think about that.
	q := datastore.NewQuery("Language").Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	q = q.Limit(10)
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
		l.Code = k.StringID()
		ls[l.Code] = l
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
	if s.R.Method == "POST" {
		l.Created = time.Now()
	}
	l.LastModified = time.Now()
	k := datastore.NewKey(s.Ctx, "Language", l.Code, 0, nil)
	k, err := datastore.Put(s.Ctx, k, l)
	return l, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities with the given limit.
func PutAndGetMulti(s *session.Session, c datastore.Cursor, l *Language) (Languages,
	datastore.Cursor, error) {
	ls := make(Languages)
	_, err := Put(s, l)
	if err != nil {
		return nil, c, err
	}
	ls, c, err = GetMulti(s, c)
	return ls, c, err
}

// GetCount "Exported functions should have a comment"
func GetCount(s *session.Session) (c int, err error) {
	q := datastore.NewQuery("Language")
	c, err = q.Count(s.Ctx)
	return
}

// Delete removes the language with the provided language code and returns an error.
func Delete(s *session.Session, langCode string) error {
	k := datastore.NewKey(s.Ctx, "Language", langCode, 0, nil)
	return datastore.Delete(s.Ctx, k)
}
