package entity

// Pass is extension for basic record type with login and password string types.
type Pass struct {
	ID   RecordID
	Login StringField
	Password StringField
}

func NewPass(key string, id RecordID, login, password StringField) *Pass {
	p := &Pass{
		ID:   id,
		Login: login,
		Password: password,
	}
	p.Login.Encrypt(key)
	p.Password.Encrypt(key)

	return p
}
