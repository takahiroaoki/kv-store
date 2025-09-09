package util

func PadStart(str string, length int, padStr string) string {
	padding := ""
	for len(padding)+len(str) <= length-len(padStr) {
		padding += padStr
	}
	padding += padStr[:length-len(padding)-len(str)]
	return padding + str
}
