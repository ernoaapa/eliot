package utils

// MergeLists merge two or more lists without duplicating items
func MergeLists(lists ...[]string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for _, list := range lists {
		for _, item := range list {
			if encountered[item] == true {
				// Do not add duplicate.
			} else {
				// Record this element as an encountered element.
				encountered[item] = true
				// Append to result slice.
				result = append(result, item)
			}
		}
	}

	return result
}

// RotateL rotates the byte array to left by one
func RotateL(a *[]string) {
	RotateLBy(a, 1)
}

// RotateLBy rotates the string array to left by given steps
func RotateLBy(a *[]string, i int) {
	x, b := (*a)[:i], (*a)[i:]
	*a = append(b, x...)
}

// RotateR rotates the string array to right by one
func RotateR(a *[]string) {
	RotateRBy(a, 1)
}

// RotateRBy rotates the string array to right by given steps
func RotateRBy(a *[]string, i int) {
	x, b := (*a)[:(len(*a)-i)], (*a)[(len(*a)-i):]
	*a = append(b, x...)
}

// GetFirst return first element from the list or default if list is empty
func GetFirst(l []string, d string) string {
	if len(l) == 0 {
		return d
	}
	return l[0]
}
