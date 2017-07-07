/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package content

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"strings"
	"time"
)

// Errors
var (
	ErrFindContent = errors.New("Error while getting content.")
	// ErrPutContent  = errors.New("Error while putting page into the datastore.")
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
func Put(ctx context.Context, c *Content) (*Content, error) {
	c.Created = time.Now()
	p.LastModified = time.Now()
	keyName := "getSha(ID)"
	k := datastore.NewKey(ctx, "Content", keyName, 0, nil)
	k, err := datastore.Put(ctx, k, p)
	p.ID = k.StringID()
	return p, err
}

func GetMulti(ctx context.Context, c datastore.Cursor) (Contents, datastore.Cursor, error) {
	ps := make(Contents)
	q := datastore.NewQuery("Content").Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	q = q.Limit(10)
	for it := q.Run(ctx); ; {
		p := new(Content)
		k, err := it.Next(p)
		if err == datastore.Done {
			c, err = it.Cursor()
			return ps, c, err
		}
		if err != nil {
			err = ErrFindContent
			return nil, c, err
		}
		p.ID = k.StringID()
		ps[p.ID] = p
	}
}
