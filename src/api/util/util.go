package util

import "unsafe"

func Int16ToInt64(bs []int16) (in int64) {
	l := len(bs)
	if l > 4 || l == 0 {
		return 0
	}

	pi := (*[4]int16)(unsafe.Pointer(&in))
	if IsBigEndian() {
		for i := range bs {
			pi[i] = bs[l-i-1]
		}
		return
	}

	for i := range bs {
		pi[3-i] = bs[l-i-1]
	}
	return
}