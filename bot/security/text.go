package security

import "strings"

func RemoveDangerousMentions(s string) string {
	return strings.ReplaceAll(s, "@everyone", "")
}
