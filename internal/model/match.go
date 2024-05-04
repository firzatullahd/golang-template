package model

type (
	InputUpdateMatch struct {
		MatchId  uint64
		Approval bool
	}

	FilterFindMatch struct {
		CatId           []uint64
		UserId          *uint64
		ID              []uint64
		Approval        *bool
		PendingApproval bool
	}

	CreateMatchRequest struct {
		UserId     uint64 `json:"-"`
		CatId      uint64 `json:"catId"`
		MatchCatId uint64 `json:"matchCatId"`
		Message    string `json:"message"`
	}

	FindMatchResponse struct {
		ID             string          `json:"id"`
		Message        string          `json:"message"`
		CreatedAt      string          `json:"createdAt"`
		UserCatDetail  FindCatResponse `json:"userCatDetail"`
		MatchCatDetail FindCatResponse `json:"matchCatDetail"`
		IssuerDetail   IssuerDetail    `json:"issuedBy"`
	}

	IssuerDetail struct {
		Email     string `json:"email"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
	}

	CatDetail struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Race string `json:"race"`
	}
)
