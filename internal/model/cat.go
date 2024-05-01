package model

type FilterFindCat struct {
	Limit      int
	Offset     int
	ID         int
	Sex        string
	Race       string
	HasMatched *bool
	Age        int
	Owned      bool
	SearchName string
	UserID     uint64
}
