package entity

// Text is basic type for all secret data.
type Text struct {
	ID   RecordID
	Text StringField
}

func NewText(key string, id RecordID, text StringField) *Text {
	t := &Text{
		ID:   id,
		Text: text,
	}
	t.Text.Encrypt(key)

	return t
}
