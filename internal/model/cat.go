package model

import "github.com/firzatullahd/cats-social-api/internal/entity"

type FilterFindCat struct {
	Limit      int
	Offset     int
	ID         int
	Sex        string
	Race       string
	HasMatched *bool
	Age        int
	SearchName string
	UserID     uint64
}

type (
	CreateCatRequest struct {
		Name        string          `json:"name"`
		SexStr      string          `json:"sex"`
		Sex         entity.SexEnum  `json:"-"`
		RaceStr     string          `json:"race"`
		Race        entity.RaceEnum `json:"-"`
		ImageUrls   []string        `json:"imageUrls"`
		AgeInMonth  uint64          `json:"ageInMonth"`
		Description string          `json:"description"`
	}
	CreateCatResponse struct {
		ID        uint64 `json:"id"`
		CreatedAt string `json:"createdAt"`
	}

	FindCatResponse struct {
		Name        string          `json:"name"`
		SexStr      string          `json:"sex"`
		Sex         entity.SexEnum  `json:"-"`
		RaceStr     string          `json:"race"`
		Race        entity.RaceEnum `json:"-"`
		ImageUrls   []string        `json:"imageUrls"`
		AgeInMonth  uint64          `json:"ageInMonth"`
		Description string          `json:"description"`
	}
)
