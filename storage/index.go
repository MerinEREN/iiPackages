/*
Package storage "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package storage

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/MerinEREN/iiPackages/crypto"
	"github.com/MerinEREN/iiPackages/session"
	"io"
	// "log"
	"mime/multipart"
	"strings"
)

const (
	gcsBucket = "inceis-1319.appspot.com"
)

// UploadFile "Exported functions should have a comment"
func UploadFile(s *session.Session, mpf multipart.File, hdr *multipart.FileHeader) (
	mLink string, err error) {
	var pfix, ext, name string
	pfix, ext, err = fileFilter(hdr)
	if err != nil {
		return
	}
	name, err = crypto.GetSha(mpf)
	if err != nil {
		return
	}
	mpf.Seek(0, 0)
	name = pfix + name + "." + ext
	mLink, err = putFile(s, mpf, name)
	return
}

// This is a weak filter, make it stronger.
func fileFilter(hdr *multipart.FileHeader) (pfix, ext string, err error) {
	ext = hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:]
	// Allow only image and video file formats.
	switch ext {
	case "jpeg", "png", "jpg", "bmp", "gif":
		pfix = "img/"
		return
	case "avi", "mov", "flv", "mpg", "mp4", "mkv":
		pfix = "video/"
		return
	default:
		err = fmt.Errorf("%s is an invalid format type to upload. Allowed types "+
			"are: jpeg, png, jpg, bmp, gif, avi, mov, flv, mpg, mp4 and mkv.",
			ext)
		return
	}
}

func putFile(s *session.Session, mpf multipart.File, name string) (
	link string, err error) {
	client := new(storage.Client)
	client, err = storage.NewClient(s.Ctx)
	if err != nil {
		return
	}
	defer client.Close()
	bucket := client.Bucket(gcsBucket)
	object := bucket.Object(name)
	w := object.NewWriter(s.Ctx)
	if err != nil {
		return
	}
	// ACLRule initialization is for make source kode able to read objects from the
	// Google Cloud Storage.
	w.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	_, err = io.Copy(w, mpf)
	if err != nil {
		return
	}
	err = w.Close()
	attrs := new(storage.ObjectAttrs)
	attrs, err = object.Attrs(s.Ctx)
	if err != nil {
		return
	}
	link = attrs.MediaLink
	return
}

// GetFile "Close ReadCloser where you call that function (defer rdr.Close())."
func GetFile(s *session.Session, name string) (io.ReadCloser, error) {
	client, err := storage.NewClient(s.Ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.Bucket(gcsBucket).Object(name).NewReader(s.Ctx)
}
