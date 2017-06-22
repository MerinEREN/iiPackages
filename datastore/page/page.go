/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package page

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"strings"
	"time"
)

// Errors
var (
	ErrFindPage = errors.New("Error while getting page.")
	// ErrPutPage  = errors.New("Error while putting page into the datastore.")
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
func Put(ctx context.Context, p *Page) (*Page, error) {
	p.Title = strings.TrimSpace(p.Title)
	// p.Path = strings.Split(p.Path, "/")[strings.LastIndex(p.Path, "/")+1]
	p.Created = time.Now()
	p.LastModified = time.Now()
	keyName := strings.Replace(p.Title, " ", "", -1)
	k := datastore.NewKey(ctx, "Page", keyName, 0, nil)
	k, err := datastore.Put(ctx, k, p)
	p.ID = k.StringID()
	return p, err
}

func GetMulti(ctx context.Context, c datastore.Cursor) (Pages, datastore.Cursor, error) {
	ps := make(Pages)
	q := datastore.NewQuery("Page").Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	q = q.Limit(10)
	for it := q.Run(ctx); ; {
		p := new(Page)
		k, err := it.Next(p)
		if err == datastore.Done {
			c, err = it.Cursor()
			return ps, c, err
		}
		if err != nil {
			err = ErrFindPage
			return nil, c, err
		}
		p.ID = k.StringID()
		ps[p.ID] = p
	}
}
