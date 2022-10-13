package model

type Session struct {
	UserID int32  `gorm:"not null"`
	Token  string `gorm:"not null"`
}
