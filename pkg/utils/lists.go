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
