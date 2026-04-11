package gox

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}
