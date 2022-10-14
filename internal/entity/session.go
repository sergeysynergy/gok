package entity

import "github.com/google/uuid"

// Session is pair values to authorize logged-in user.
type Session struct {
	UserID UserID
	Token  string
}

func NewSession(usrID UserID) *Session {
	s := &Session{
		UserID: usrID,
	}
	s.genToken()

	return s
}

// genToken generate new session token.
func (s *Session) genToken() {
	s.Token = uuid.New().String()
}
