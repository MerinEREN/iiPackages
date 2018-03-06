package pageContent

// PageContent ContentID is content's key's intID as string
// and PageID is page's key's stringID.
// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type PageContent struct {
	ContentID string
	PageID    string
}

// PageContents is a []*PageContent
type PageContents []*PageContent
