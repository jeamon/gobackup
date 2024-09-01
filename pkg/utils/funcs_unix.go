//go:build !windows
// +build !windows

package utils

// fixColonCharacter does nothing non-windows machines.
func FixColonCharacter(s string) string {
	return s
}
