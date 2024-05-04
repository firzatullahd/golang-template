package model

type (
	FilterFindCat struct {
		Limit      int
		Offset     int
		ID         *uint64
		Sex        *string
		Race       *string
		HasMatched *bool
		Age        *int
		SearchName *string
		UserID     *uint64
	}

	InputUpdateCat struct {
		ID          uint64
		UserID      uint64
		Name        *string
		Sex         *string
		Race        *string
		ImageUrls   []string
		Age         *uint64
		Description *string
	}
)

type (
	CreateCatRequest struct {
		UserID      uint64   `json:"-"`
		Name        string   `json:"name"`
		Sex         string   `json:"sex"`
		Race        string   `json:"race"`
		ImageUrls   []string `json:"imageUrls"`
		AgeInMonth  uint64   `json:"ageInMonth"`
		Description string   `json:"description"`
	}

	UpdateCatRequest struct {
		ID          uint64   `json:"-"`
		UserID      uint64   `json:"-"`
		Name        *string  `json:"name"`
		Sex         *string  `json:"sex"`
		Race        *string  `json:"race"`
		ImageUrls   []string `json:"imageUrls"`
		AgeInMonth  *uint64  `json:"ageInMonth"`
		Description *string  `json:"description"`
	}

	CreateCatResponse struct {
		ID        string `json:"id"`
		CreatedAt string `json:"createdAt"`
	}

	FindCatResponse struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Sex         string   `json:"sex"`
		Race        string   `json:"race"`
		ImageUrls   []string `json:"imageUrls"`
		AgeInMonth  uint64   `json:"ageInMonth"`
		Description string   `json:"description"`
		HasMatched  bool     `json:"hasMatched"`
		CreatedAt   string   `json:"createdAt"`
	}
)
