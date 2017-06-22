/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package storage

import (
	"cloud.google.com/go/storage"
	"github.com/MerinEREN/iiPackages/session"
	"github.com/MerinEREN/iiPackages/crypto"
	"mime/multipart"
	"strings"
	"io/ioutil"
	"errors"
)

const (
	ErrInvalidFileType = errors.New("Invalid file type.")
)

var (
	app        = "inceis-1319"
	bucketName = "inceis-1319.appspot.com"
)

func UploadFile(s *session.Session, mpf multipart.File, hdr *multipart.FileHeader) (
	string, error) {
	client, err := storage.NewClient(s.ctx)
	defer client.Close()
	bucket := client.Bucket(bucketName)
	ext, err := fileFilter(hdr)
	if err != nil {
		return nil, err
	}
	name, err := crypto.GetSha(mpf) + '.' + ext
	if err != nil {
		return nil, err
	}
	object := bucket.Object(name)
}


func putFile(s *session.Session, f *multipart.File) error {
}

func fileFilter(hdr *multipart.FileHeader) (ext string, err error) {
	ext = hdr[strings.LastIndex(hdr, ".")+1:]
	// Allow only image and video file formats.
	switch ext {
	case 'jpeg', 'png', 'jpg', 'bmp', 'gif', 'avi', 'mov', 'flv', 'mpg', 'mp4', 'mkv':
		return ext, nil
	default:
		return nil, ErrInvalidFileType
	}
}
