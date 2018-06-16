/*
Package tag "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package tag

import (
	"errors"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindTag = errors.New("Error while getting tag")
)

// GetMulti returns limited entitity from the given cursor.
// If limit is nil default limit will be used.
func GetMulti(s *session.Session, c datastore.Cursor, limit interface{}) (Tags, datastore.Cursor, error) {
	ts := make(Tags)
	q := datastore.NewQuery("Tag").Project("Name").Order("-Created")
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
		t := new(Tag)
		k, err := it.Next(t)
		if err == datastore.Done {
			c, err = it.Cursor()
			return ts, c, err
		}
		if err != nil {
			err = ErrFindTag
			return nil, c, err
		}
		t.ID = k.Encode()
		ts[t.ID] = t
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
func Put(s *session.Session, t *Tag) (*Tag, error) {
	k := datastore.NewIncompleteKey(s.Ctx, "Tag", nil)
	var err error
	t.Created = time.Now()
	k, err = datastore.Put(s.Ctx, k, t)
	t.ID = k.Encode()
	return t, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities by the given limit.
func PutAndGetMulti(s *session.Session, c datastore.Cursor, t *Tag) (Tags,
	datastore.Cursor, error) {
	ts := make(Tags)
	tNew := new(Tag)
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		tNew, err1 = Put(s, t)
		if err1 != nil {
			return
		}
		ts, c, err1 = GetMulti(s, c, 19)
		return
	}, nil)
	ts[tNew.ID] = tNew
	return ts, c, err
}

// Delete removes the entity by the provided encoded entity key and returns an error.
func Delete(s *session.Session, ek string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	return datastore.Delete(s.Ctx, k)
}

// DeleteMulti removes the entitys by the provided encoded entity key slice
// and returns an error.
/* func DeleteMulti(s *session.Session, ekx []string) error {
	var kx []*datastore.Key
	for _, v := range ekx {
		k, err := datastore.DecodeKey(v)
		if err != nil {
			return err
		}
		kx = append(kx, k)
	}
	return datastore.DeleteMulti(s.Ctx, kx)
}*/
