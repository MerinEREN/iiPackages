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
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// GetAll returns all the entities from the begining of the kind.
func GetAll(ctx context.Context) (Languages, error) {
	var lx []*Language
	q := datastore.NewQuery("Language")
	q = q.
		Project("ContextID", "Link").
		Order("-Created")
	kx, err := q.GetAll(ctx, lx)
	if err != nil {
		return nil, err
	}
	ls := make(Languages)
	for i, v := range kx {
		lx[i].ID = v.StringID()
		ls[v.StringID()] = lx[i]
	}
	return ls, err
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

// PutAndGetAll is a transaction which puts the posted item first
// and then gets entities by the given limit.
func PutAndGetAll(s *session.Session, l *Language) (Languages, error) {
	ls := make(Languages)
	lNew := new(Language)
	// USAGE "s" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!!!!!!!!!
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		lNew, err1 = Put(s, l)
		if err1 != nil {
			return
		}
		ls, err1 = GetAll(s.Ctx)
		return
	}, nil)
	ls[lNew.ID] = lNew
	return ls, err
}

// GetCount returns language count and an error.
/* func GetCount(ctx context.Context) (c int, err error) {
	q := datastore.NewQuery("Language")
	c, err = q.Count(ctx)
	return
} */
