package buffer

import (
	"strings"
	"testing"
)

func TestNewPool(t *testing.T) {
	p := NewPool(10)
	b := p.Get()
	if b == nil {
		t.Fatalf("TestNewPool falied")
	}
	b.WriteString("a")
	if b.String() != "a" {
		t.Fatalf("TestNewPool falied")
	}
	p.Put(b)
	b = p.Get()
	if b == nil || b.Len() != 0 {
		t.Fatalf("TestNewPool falied")
	}
}

func BenchmarkNewPool(b *testing.B) {
	p := NewPool(4 * 10)
	s := strings.Repeat("a", 4*1024)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := p.Get()
			buf.WriteString(s)
			_ = buf.String()
			p.Put(buf)
		}
	})
	b.ReportAllocs()
	// 2000000	       872 ns/op	    4098 B/op	       1 allocs/op
}
