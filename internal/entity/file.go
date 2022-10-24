package entity

// File is extension for basic record type to store binary data.
type File struct {
	ID   RecordID
	File []byte
}

func NewFile(key string, id RecordID, file []byte) *File {
	f := &File{
		ID:   id,
		File: file,
	}

	return f
}
