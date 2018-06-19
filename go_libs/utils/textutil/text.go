package textutil
import "unicode"
func IsChineseChar(str string) bool {
	han := unicode.Scripts["Han"]
	for _, r := range str {
		if unicode.Is(han, r) {
			return true
		}
	}
	return false
}
