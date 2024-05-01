package entity

import "time"

type Match struct {
	ID          string `db:"id"`
	UserID      int    `db:"user_id"`
	CatID       int    `db:"cat_id"`
	MatchUserID int    `db:"match_user_id"`
	MatchCatID  int    `db:"match_cat_id"`

	IsApproved bool `db:"is_approved"`
	IsDeclined bool `db:"is_declined"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
