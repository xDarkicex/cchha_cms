package controllers

import (
	"encoding/base64"

	"github.com/xDarkicex/hasher"
)

// Avalible Hashers
var (
	MD5    = hasher.NewHasher().MD5()
	SHA1   = hasher.NewHasher().SHA1()
	SHA256 = hasher.NewHasher().SHA256()
	SHA512 = hasher.NewHasher().SHA512()
)

// DecodeBASE64 taked encode in base64 URL encoding returns []byte
// for our use sent Decrypt uuid
func DecodeBASE64(encoded string) (decoded []byte, err error) {
	return base64.URLEncoding.DecodeString(encoded)
}

// EncodeBASE64 takes []bytes returns string
func EncodeBASE64(plainstring []byte) string {
	return base64.URLEncoding.EncodeToString(plainstring)
}
