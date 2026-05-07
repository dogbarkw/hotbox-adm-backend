package util

import "github.com/speps/go-hashids/v2"

const (
	HashSalt = "7cd30abtzna8"
)

func EncodeInt64(id int64) (e string) {
	hd := hashids.NewData()
	hd.Salt = HashSalt
	h, err := hashids.NewWithData(hd)
	if nil != err {
		return ""
	}
	e, err = h.EncodeInt64([]int64{id})
	if nil != err {
		return ""
	}
	return e
}
