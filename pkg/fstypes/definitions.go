package fstypes

type Type string

const (
	FILE        Type = "FILE"
	DIR         Type = "DIRECTORY"
	SYMLINK     Type = "SYMLINK"
	UNSUPPORTED Type = "UNSUPPORTED"
)
