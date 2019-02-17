package util

// StringSlice ...
type StringSlice []string

// Has returns whether the str is in the slice.
func (s StringSlice) Has(str string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == str {
			return true
		}
	}
	return false
}

// Add adds the str to the slice.
func (s StringSlice) Add(str string) ([]string, bool) {
	if s.Has(str) {
		return s, false
	}
	return append(s, str), true
}

// Remove remove the str from the slice.
func (s StringSlice) Remove(str string) ([]string, bool) {
	offset := 0
	for i := 0; i < len(s); i++ {
		if s[i] != str {
			s[offset] = s[i]
			offset++
		}
	}
	if offset < len(s) {
		return s[0:offset], true
	}
	return s, false
}
