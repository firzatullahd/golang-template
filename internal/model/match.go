package model

type (
	InputUpdateMatch struct {
		MatchId     uint64
		MatchUserId uint64
		MatchCatId  uint64
		Approval    bool
	}

	FilterFindMatch struct {
		CatId  *uint64
		UserId *uint64
	}
)
