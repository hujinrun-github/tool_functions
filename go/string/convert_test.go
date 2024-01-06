package convert

import "testing"

func BenchmarkStr2Bytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		str2Bytes("34dfafdas32132")
	}
}

func BenchmarkOriginStr2Bytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		originStr2Bytes("34dfafdas32132")
	}
}

func originStr2Bytes(s string) []byte {
	return []byte(s)
}
