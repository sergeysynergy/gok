package entity

type BranchID uint32

type Branch struct {
	ID     BranchID
	UserID UserID
	Name   string
	Head   uint64
}
