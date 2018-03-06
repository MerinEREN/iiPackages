/*
Package datastore has utility functions to use with queries.
*/
package datastore

import (
	"google.golang.org/appengine/datastore"
)

// FilterMulti applays multiple filter operations on same property.
func FilterMulti(q *datastore.Query, str string, slc interface{}) *datastore.Query {
	switch s := slc.(type) {
	case []*datastore.Key:
		for _, v := range s {
			q = q.Filter(str, v)
		}
	}
	return q
}
