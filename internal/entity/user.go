package entity

import "github.com/google/uuid"

// UserID type provides more easy user id type changes.
type UserID int32

// User defines user structure, where:
// ID uniq user key;
// Login uniq username;
// Key used to encrypt user records data.
type User struct {
	ID    UserID
	Login string
	Key   string
}

// SignedUser store data to give response for successfully signed-in users.
type SignedUser struct {
	Token string
	Key   string
}

// NewSignedUser creates new SignedUser with key.
func NewSignedUser(token string) *SignedUser {
	u := &SignedUser{
		Token: token,
	}
	u.genKey()

	return u
}

// genKey generates an individual key to use for user records encryption.
func (s *SignedUser) genKey() {
	s.Key = uuid.New().String()
}

type CLIUser struct {
	Login string
	Home  string
}
