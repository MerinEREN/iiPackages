/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

// Errors
var (
	ErrInvalidFileType = errors.New("Invalid file type to encrypt.")
)

func GetMAC(i interface{}) (string, error) {
	mac := hmac.New(sha256.New, []byte("someKey"))
	switch v := i.(type) {
	case string:
		io.WriteString(mac, v)
	case io.Reader:
		io.Copy(mac, v)
	default:
		return "", ErrInvalidFileType
	}
	return fmt.Sprintf("%x", mac.Sum(nil)), nil
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func CheckMAC(message interface{}, messageMAC string) bool {
	// mac := hmac.New(sha256.New, key)
	// mac.Write(message)
	// expectedMAC := mac.Sum(nil)
	// return hmac.Equal(messageMAC, expectedMAC)
	macString, _ := GetMAC(message)
	expectedMAC := []byte(macString)
	return hmac.Equal([]byte(messageMAC), expectedMAC)
}

func GetSha(i interface{}) (string, error) {
	h := sha256.New()
	switch v := i.(type) {
	case string:
		io.WriteString(h, v)
	case io.Reader:
		io.Copy(h, v)
	case *os.File:
		io.Copy(h, v)
	default:
		return "", ErrInvalidFileType
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
