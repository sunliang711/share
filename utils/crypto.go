package utils

import (
	"crypto/md5"
)

func Md5(bs []byte) []byte {
	h := md5.New()
	h.Write(bs)
	return h.Sum(nil)
}
