/*
Package page "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package page

import (
	"errors"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"strings"
	"time"
)

// Errors
var (
	ErrFindPage = errors.New("error while getting page")
	// ErrPutPage  = errors.New("Error while putting page into the datastore.")
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
func Put(s *session.Session, p *Page) (Pages, error) {
	ps := make(Pages)
	k := new(datastore.Key)
	p.Title = strings.TrimSpace(p.Title)
	if s.R.Method == "POST" {
		p.Created = time.Now()
		keyName := strings.Replace(p.Title, " ", "", -1)
		k = datastore.NewKey(s.Ctx, "Page", keyName, 0, nil)
		p.ID = k.StringID()
	} else if s.R.Method == "PUT" {
		k = datastore.NewKey(s.Ctx, "Page", p.ID, 0, nil)
	}
	p.LastModified = time.Now()
	k, err := datastore.Put(s.Ctx, k, p)
	ps[p.ID] = p
	return ps, err
}

// GetMulti "Exported functions should have a comment"
func GetMulti(s *session.Session, c datastore.Cursor) (Pages, datastore.Cursor, error) {
	ps := make(Pages)
	// Maybe -LastModied should be the order ctireia if consider UX, think about that.
	q := datastore.NewQuery("Page").Project("Title", "Link").Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	q = q.Limit(10)
	for it := q.Run(s.Ctx); ; {
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

// Get returns the page with provided keyName and an error.
func Get(s *session.Session, keyName string) (Pages, error) {
	ps := make(Pages)
	k := datastore.NewKey(s.Ctx, "Page", keyName, 0, nil)
	p := new(Page)
	err := datastore.Get(s.Ctx, k, p)
	p.ID = k.StringID()
	ps[p.ID] = p
	return ps, err
}

// Delete removes the page with the provided page id and returns an error.
func Delete(s *session.Session, pageID string) error {
	k := datastore.NewKey(s.Ctx, "Page", pageID, 0, nil)
	return datastore.Delete(s.Ctx, k)
}
