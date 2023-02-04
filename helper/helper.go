package helper

func RemoveDuplicateString(oldList []string) []string {
	keys := make(map[string]bool)
	var newList []string

	for _, entry := range oldList {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			newList = append(newList, entry)
		}
	}
	return newList
}
