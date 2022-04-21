package stringandbytes

import (
	"reflect"
	"unsafe"
)

// String2Bytes 将 string 强转换为 []byte ，性能高于标准转换 []byte(string),
//但可能造成危险(如对转换后的 []byte 操作会导致对原本不能被修改的 string 被修改,并且不能被 defer + recover 捕获错误)
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
