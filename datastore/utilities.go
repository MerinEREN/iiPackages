/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package datastore

import (
	"google.golang.org/appengine/datastore"
)

func FilterMulti(q *datastore.Query, str string, slc interface{}) *datastore.Query {
	switch s := slc.(type) {
	case []*datastore.Key:
		for _, v := range s {
			q = q.Filter(str, v)
		}
	}
	return q
}
