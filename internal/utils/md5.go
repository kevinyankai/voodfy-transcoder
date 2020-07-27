package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

func CreateHash(id string) string {
	key := fmt.Sprintf("%s%s", id, RandSeq(256))
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
