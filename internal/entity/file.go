package entity

// File is extension for basic record type to store binary data.
type File struct {
	ID   RecordID
	File StringField
}

func NewFile(key string, id RecordID, file StringField) *File {
	f := &File{
		ID:   id,
		File: file,
	}
	f.File.Encrypt(key)

	return f
}
