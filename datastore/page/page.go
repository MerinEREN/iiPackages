/*
Package page "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package page

import (
	"errors"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindPage = errors.New("error while getting page")
)

/*
GetMulti returns corresponding pages if page keys provided
Otherwise returns limited entitity from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(s *session.Session, c datastore.Cursor, limit, kx interface{}) (Pages, datastore.Cursor, error) {
	ps := make(Pages)
	if kx, ok := kx.([]*datastore.Key); ok {
		var px []*Page
		err := datastore.GetMulti(s.Ctx, kx, px)
		if err != nil {
			return nil, c, err
		}
		for i, v := range kx {
			px[i].ID = v.StringID()
			ps[v.StringID()] = px[i]
		}
		return ps, c, err
	}
	// Maybe -LastModied should be the order ctireia if consider UX, think about that.
	q := datastore.NewQuery("Page").
		Project("Title", "Link").
		Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	if limit != nil {
		l := limit.(int)
		q = q.Limit(l)
	} else {
		q = q.Limit(10)
	}
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
func Put(s *session.Session, p *Page) (*Page, error) {
	k := datastore.NewKey(s.Ctx, "Page", p.ID, 0, nil)
	var err error
	p.LastModified = time.Now()
	if s.R.Method == "POST" {
		p.Created = time.Now()
		_, err = datastore.Put(s.Ctx, k, p)
	} else if s.R.Method == "PUT" {
		// For "PUT" request
		tempP := new(Page)
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			if err1 = datastore.Get(ctx, k, tempP); err1 != nil {
				return
			}
			p.Created = tempP.Created
			_, err1 = datastore.Put(ctx, k, p)
			return
		}, nil)
	}
	return p, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities with the given limit.
func PutAndGetMulti(s *session.Session, c datastore.Cursor, p *Page) (Pages,
	datastore.Cursor, error) {
	ps := make(Pages)
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) error {
		p, err1 := Put(s, p)
		if err1 != nil {
			return err1
		}
		ps, c, err1 = GetMulti(s, c, 9, nil)
		ps[p.ID] = p
		return err1
	}, nil)
	return ps, c, err
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
