package pageContent

import (
	"google.golang.org/appengine/datastore"
)

// PageContent datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!
// Page key is the parent key.
type PageContent struct {
	ContentKey *datastore.Key
}

// PageContents is a []*PageContent
type PageContents []*PageContent
