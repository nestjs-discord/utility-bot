package security

import "strings"

func RemoveDangerousMentions(s string) string {
	return strings.Replace(s, "@everyone", "", -1)
}
