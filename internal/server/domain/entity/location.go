package entity

type LocationInd int32

type Location struct {
	Ind         LocationInd
	Caption     string
	Description *string
}

type ByInd map[LocationInd]*Location

type SortByInd []*Location
