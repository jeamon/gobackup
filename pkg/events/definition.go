package events

import "github.com/jeamon/gobackup/pkg/fstypes"

// Change represents the object to be processed by backup workers.
type Change struct {
	Path  string
	Ops   Event
	Type  fstypes.Type
	Error error
}

// EventChan is a channel of change events.
type Queue chan *Change

// Event is custom type to restrict possible event values.
type Event string

const (
	// This event denotes a file creation.
	CREATE Event = "CREATE"
	// This event denotes a file modification.
	MODIFY Event = "MODIFY"
	// This event denotes a file deletion.
	DELETE Event = "DELETE"
	// This event denotes a file renaming.
	RENAME Event = "RENAME"
	// This event denotes a file attributes change.
	ATTRIBUTE Event = "PERM"
	// This event denotes a folder creation.
	WATCH Event = "WATCH"
	// Request an immediate or scheduled deletion.
	RDELETE Event = "RREMOVE"
	// This event means nothing change on filesystem.
	NOCHANGE Event = "NOCHANGE"
)
