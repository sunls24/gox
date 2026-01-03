package gox

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand/v2"
	"unsafe"
)

func MD5(str string) string {
	sum := md5.Sum(Str2Bytes(str))
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum[:])
	return Bytes2Str(dst)
}

func RandStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return Bytes2Str(b)
}

func Str2Bytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func Bytes2Str(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
