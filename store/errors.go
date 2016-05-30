package store

// FileNotFoundError is returned by a store if a file does not exist in
// the store.
type FileNotFoundError string

func (e FileNotFoundError) Error() string {
	return string(e)
}

// ExpirationInPastError is returned by a store's Store function if a File's Metadata Expire field
// holds a date in the past.
type ExpirationInPastError string

func (e ExpirationInPastError) Error() string {
	return string(e)
}

// CorruptedDataError should be returned if a store became unsuable.
type CorruptedDataError string

func (e CorruptedDataError) Error() string {
	return string(e)
}
