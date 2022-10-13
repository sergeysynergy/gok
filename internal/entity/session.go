package entity

type SessionID int32

type Session struct {
	UserID SessionID
	Token  string
}
