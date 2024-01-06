package convert

import (
	"reflect"
	"unsafe"
)

func str2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
