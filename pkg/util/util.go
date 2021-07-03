package util

// Exist check element has exist in elements or not
func Exist(element string, elements []*string) bool {
	b := false
	for _, i := range elements {
		if *i == element {
			b = true
			return b
		}
	}
	return b
}
