package entity

import (
	"fmt"
	"strings"
	"time"
)

type Cat struct {
	ID     uint64 `db:"id"`
	UserID uint64 `db:"user_id"`
	Name   string `db:"name"`

	Sex         SexEnum  `db:"sex"`
	Race        RaceEnum `db:"race"`
	ImageUrls   []string `db:"image_urls"`
	Age         uint64   `db:"age"`
	Description string   `db:"description"`
	HasMatched  bool     `db:"has_matched"`

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

func StringToSex(str string) (result SexEnum, err error) {
	switch strings.ToUpper(str) {
	case "MALE":
		result = SexMale
	case "FEMALE":
		result = SexFemale
	default:
		err = fmt.Errorf("sex(%s) doesn't recognized", str)
		return
	}
	return
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

func StringToRace(str string) (result RaceEnum, err error) {
	switch strings.ToUpper(str) {
	case "Persian":
		result = RacePersian
	case "Maine Coon":
		result = RaceMaineCoon
	case "Siamese":
		result = RaceSiamese
	case "Ragdoll":
		result = RaceRagdoll
	case "RaceBengal":
		result = RaceBengal
	case "Sphynx":
		result = RaceSphynx
	case "British Shorthair":
		result = RaceBSH
	case "Abyssinian":
		result = RaceAbys
	case "Scottish Fold":
		result = RaceScottish
	case "Birman":
		result = RaceBirman
	default:
		err = fmt.Errorf("race(%s) doesn't recognized", str)
		return
	}
	return
}
