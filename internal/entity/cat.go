package entity

import "time"

// todo remove json
type Cat struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	Name   string `json:"name" db:"name"`

	ImageUrls [][]string
	Sex       sexEnum
	Race      raceEnum

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type sexEnum int

const (
	SexMale sexEnum = iota + 1
	SexFemale
)

func (e sexEnum) String() string {
	switch e {
	case SexMale:
		return "MALE"
	case SexFemale:
		return "FEMALE"
	default:
		return "UNKNOWN"
	}
}

type raceEnum int

const (
	RacePersian raceEnum = iota + 1
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

func (e raceEnum) String() string {
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
