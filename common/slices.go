package common

// IndexOf finds the position in a slice of an given value.
func IndexOf(stack []interface{}, needle interface{}) int {
	for index := 0; index < len(stack); index++ {
		if stack[index] == needle {
			return index
		}
	}
	return -1
}
