package dto

import (
	"fmt"
	"regexp"
)

var nonDigit = regexp.MustCompile(`\D`)

func formatCNPJ(s string) string {
	d := nonDigit.ReplaceAllString(s, "")
	if len(d) != 14 {
		return s
	}
	return fmt.Sprintf("%s.%s.%s/%s-%s", d[0:2], d[2:5], d[5:8], d[8:12], d[12:14])
}

func formatCPF(s string) string {
	d := nonDigit.ReplaceAllString(s, "")
	if len(d) != 11 {
		return s
	}
	return fmt.Sprintf("%s.%s.%s-%s", d[0:3], d[3:6], d[6:9], d[9:11])
}

func formatPhone(s string) string {
	d := nonDigit.ReplaceAllString(s, "")
	switch len(d) {
	case 11:
		return fmt.Sprintf("(%s) %s-%s", d[0:2], d[2:7], d[7:11])
	case 10:
		return fmt.Sprintf("(%s) %s-%s", d[0:2], d[2:6], d[6:10])
	}
	return s
}

func formatZipCode(s string) string {
	d := nonDigit.ReplaceAllString(s, "")
	if len(d) != 8 {
		return s
	}
	return fmt.Sprintf("%s-%s", d[0:5], d[5:8])
}
