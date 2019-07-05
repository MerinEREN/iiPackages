package pageContext

import (
	"google.golang.org/appengine/datastore"
)

// PageContext datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!
// Page key is the parent key.
type PageContext struct {
	ContextKey *datastore.Key
}

// PageContexts is a []*PageContext
type PageContexts []*PageContext
