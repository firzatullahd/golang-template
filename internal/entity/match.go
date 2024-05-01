package entity

import "time"

type Match struct {
	ID          uint64 `db:"id"`
	UserID      uint64 `db:"user_id"`
	CatID       uint64 `db:"cat_id"`
	MatchUserID uint64 `db:"match_user_id"`
	MatchCatID  uint64 `db:"match_cat_id"`

	IsApproved bool   `db:"is_approved"`
	IsRejected bool   `db:"is_rejected"`
	Message    string `db:"message"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
