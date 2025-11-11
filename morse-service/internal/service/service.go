package service

import (
	"errors"
	"strings"

	"sprint6/pkg/morse"
)

func isMorseLike(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	for _, r := range s {
		switch r {
		case '.', '-', ' ', '/', '\n', '\r', '\t':
			continue
		default:
			return false
		}
	}
	return strings.ContainsAny(s, ".-")
}

func ConvertAuto(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("пустые данные: нечего конвертировать")
	}
	if isMorseLike(trimmed) {
		return morse.ToText(trimmed), nil
	}
	return morse.ToMorse(trimmed), nil
}
