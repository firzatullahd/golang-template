package entity

import "time"

type Cat struct {
	ID     uint64 `db:"id"`
	UserID uint64 `db:"user_id"`
	Name   string `db:"name"`

	Sex         SexEnum    `db:"sex"`
	Race        RaceEnum   `db:"race"`
	ImageUrls   [][]string `db:"image_urls"`
	Age         uint64     `db:"age"`
	Description string     `db:"description"`
	HasMatched  bool       `db:"has_matched"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type SexEnum int

const (
	SexMale SexEnum = iota + 1
	SexFemale
)

func (e SexEnum) String() string {
	switch e {
	case SexMale:
		return "MALE"
	case SexFemale:
		return "FEMALE"
	default:
		return "UNKNOWN"
	}
}

type RaceEnum int

const (
	RacePersian RaceEnum = iota + 1
	RaceMaineCoon
	RaceSiamese
	RaceRagdoll
	RaceBengal
	RaceSphynx
	RaceBSH
	RaceAbys
	RaceScottish
	RaceBirman
)

func (e RaceEnum) String() string {
	switch e {
	case RacePersian:
		return "Persian"
	case RaceMaineCoon:
		return "Maine Coon"
	case RaceSiamese:
		return "Siamese"
	case RaceRagdoll:
		return "Ragdoll"
	case RaceBengal:
		return "Bengal"
	case RaceSphynx:
		return "Sphynx"
	case RaceBSH:
		return "British Shorthair"
	case RaceAbys:
		return "Abyssinian"
	case RaceScottish:
		return "Scottish Fold"
	case RaceBirman:
		return "Birman"

	default:
		return "UNKNOWN"
	}
}
