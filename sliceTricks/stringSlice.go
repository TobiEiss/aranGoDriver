package sliceTricks

// Find a element in slice
// index := session.database.Find(func(index int, value string) bool {
//		return value == term
//	})
func Find(slice []string, find func(index int, value string) bool) int {
	for index, value := range slice {
		if find(index, value) {
			return index
		}
	}
	return -1
}

// Contains check if a string is contained in a string slice
func Contains(slice []string, term string) bool {
	index := Find(slice, func(index int, value string) bool {
		return value == term
	})
	return (index >= 0)
}
