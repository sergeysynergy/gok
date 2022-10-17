package model

type Session struct {
	UserID int32  `gorm:"not null;uniqueIndex:UserIDAndToken"`
	Token  string `gorm:"not null;uniqueIndex:UserIDAndToken"`
}
