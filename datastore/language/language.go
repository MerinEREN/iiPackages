/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package language

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindLanguage = errors.New("Error while getting language.")
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
func Put(ctx context.Context, l *Language) (*Language, error) {
	l.Created = time.Now()
	l.LastModified = time.Now()
	k := datastore.NewKey(ctx, "Language", l.Code, 0, nil)
	k, err := datastore.Put(ctx, k, l)
	return l, err
}

func GetMulti(ctx context.Context, c datastore.Cursor) (Languages, datastore.Cursor, error) {
	ls := make(Languages)
	q := datastore.NewQuery("Language").Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	q = q.Limit(10)
	for it := q.Run(ctx); ; {
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

func GetCount(ctx context.Context) (c int, err error) {
	q := datastore.NewQuery("Language")
	c, err = q.Count(ctx)
	return
}
