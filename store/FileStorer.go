package store

import (
	"io"
	"time"

	"github.com/pborman/uuid"
)

// MetaData describes a file.
type MetaData struct {
	// Name is the name of the file as it should be returned.
	Name string

	// ID is the unique indetifier by which a fille will be looked up in the
	// underlying store
	ID uuid.UUID

	// Expires gives the time after which a file is no longer valid.
	// The file may be removed from the database at any point in time after the expiry date
	Expires time.Time

	// Public is a flag which denotes wether a file will be displayed to unauthenticated users.
	Public bool
}

// File consists of metadata describing it and the files actual content
type File struct {
	Metadata MetaData
	Content  io.ReadWriter
}

// The FileStorer interface describes the behaviour every file store has to implement.
type FileStorer interface {

	// IsPublic returns true if the file is free to be displayed publicly
	// and false if not. Defaults to false for security reasons.
	IsPublic(uuid.UUID) bool

	// Is Expired returns true if the files expiry date is in the past, false if not.
	// Defaults to false, so that a file will still be available, even when the expiration
	// state is undefined.
	IsExpired(uuid.UUID) bool

	// Store takes an instance of File and stores it in the underlying engine.
	// It should store the files Metadata as well.
	// Returns an nil and an ExpirationInPastError if the metadata's Expire field
	// holds a date in the past.
	// If a file is suvessfully stored, it's Metadata's ID field gets updated.
	// On successfully saving the file, the files UUID and a nil error is returned.
	// Should an error occur while storing the file, a nil UUID and the respective error is returned.
	//
	// A file given to store may already have a UUID in its metadata. If no UUID exists, one is generated.
	Store(file *File) (uuid.UUID, error)

	// Get takes a UUID and retrieves the file associated with it.
	// It is the implementations duty to make sure the files metadata is populated.
	// Returns nil and a FileNotFoundError if a file does not exist in the store
	Get(uuid.UUID) (*File, error)

	// Remove takes a UUID and removes the file from the underlying storage
	Remove(uuid.UUID) error
}
