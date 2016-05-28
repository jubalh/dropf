package store

import "io"

// FileStorer defines an interface for pluggable file storage.
type FileStorer interface {
	// IsPublic sets a flag to show whether a file can be accessed by everybody or just registered users.
	IsPublic() bool

	// StoreFile stores a file. It uses buffered IO.
	StoreFile(io.Reader) error

	// GetFile gets a file by its UUID.
	GetFile(uuid uuid.UUID) io.Reader

	// DeleteFile deletes a file by its UUID.
	DeleteFile(uuid uuid.UUID) error
}
