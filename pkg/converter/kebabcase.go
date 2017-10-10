package converter

import "strings"

// KebabCaseToCamelCase converts kebab-case string to CamelCase
func KebabCaseToCamelCase(kebab string) (camelCase string) {
	isToUpper := false
	for i, runeValue := range kebab {

		if isToUpper || i == 0 {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '-' {
				isToUpper = true
			} else {
				camelCase += string(runeValue)
			}
		}
	}
	return
}
