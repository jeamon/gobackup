//go:build windows
// +build windows

package utils

// fixColonCharacter replaces all character rune equals
// to `61498` by colon. This helps fix path name encoding
// issue on windows machine.
func FixColonCharacter(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if r == 61498 {
			runes[i] = 58
			continue
		}
	}
	return string(runes)
}
