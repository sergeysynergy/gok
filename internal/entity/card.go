package entity

// Card is extension for basic record type to safely store credit card essentials.
type Card struct {
	ID      RecordID
	Number  NumberField
	Code    NumberField
	Expired StringField
	Owner   StringField
}

func NewCard(key string, id RecordID, number, code NumberField, expired, owner StringField) *Card {
	c := &Card{
		ID:      id,
		Number:  number,
		Code:    code,
		Expired: expired,
		Owner:   owner,
	}
	c.Number.Encrypt(key)
	c.Code.Encrypt(key)
	c.Expired.Encrypt(key)
	c.Owner.Encrypt(key)

	return c
}
