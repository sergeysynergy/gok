package entity

type UserID int32

type User struct {
	ID    UserID
	Login string
	Key   string
}

// KeyGen generates an individual key to use for user records encryption.
func (u *User) KeyGen() {
	// TODO: Add key generator
	u.Key = "CodeKeyForAllUsersSoFar"
}

type UsersList struct {
	Users []*User
}

type SignedUser struct {
	Token string
	Key   string
}
