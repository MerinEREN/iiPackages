/*
Package content "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package content

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindContent = errors.New("error while getting content")
	// ErrPutContent  = errors.New("Error while putting page into the datastore.")
)

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
func Put(ctx context.Context, c *Content) (*Content, error) {
	c.Created = time.Now()
	c.LastModified = time.Now()
	keyName := "getSha(ID)"
	k := datastore.NewKey(ctx, "Content", keyName, 0, nil)
	k, err := datastore.Put(ctx, k, c)
	c.ID = k.StringID()
	return c, err
}

// GetMulti returns page contents with all languages from the cursor if available
// nor from the beginning of the kind for given limit.
func GetMulti(ctx context.Context, c datastore.Cursor) (Contents, datastore.Cursor, error) {
	cntnts := make(Contents)
	q := datastore.NewQuery("Content").Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	q = q.Limit(10)
	for it := q.Run(ctx); ; {
		cntnt := new(Content)
		k, err := it.Next(cntnt)
		if err == datastore.Done {
			c, err = it.Cursor()
			return cntnts, c, err
		}
		if err != nil {
			err = ErrFindContent
			return nil, c, err
		}
		cntnt.ID = k.StringID()
		cntnts[cntnt.ID] = cntnt
	}
}
