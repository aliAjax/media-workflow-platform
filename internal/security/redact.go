package security

import (
	"regexp"
	"strings"
)

var secretPattern = regexp.MustCompile(`(?i)(authorization|cookie|token|secret|password)=?[^\s,;]+`)

func Redact(s string) string { return secretPattern.ReplaceAllString(s, "$1=[REDACTED]") }
func RedactMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		lk := strings.ToLower(k)
		if strings.Contains(lk, "token") || strings.Contains(lk, "secret") || strings.Contains(lk, "password") || lk == "authorization" {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
