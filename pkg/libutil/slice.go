package libutil

// remove emtpy string from string slice
func RemoveEmptyString(strSlice []string) []string {
	var ret []string
	for _, str := range strSlice {
		if str != "" {
			ret = append(ret, str)
		}
	}
	return ret
}
