package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type User struct {
	ID         uint64    `db:"id"`
	Username   string    `db:"username"`
	Password   string    `db:"password"`
	Name       string    `db:"name"`
	IdCardNo   *string   `db:"id_card_no"`
	IdcardFile *string   `db:"id_card_file"`
	State      UserState `db:"state"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserState int

const (
	UserStateRegistered UserState = iota + 1
	UserStatePending
	UserStateVerified
	UserStateDeleted
)

func (s UserState) String() string {
	switch s {
	case UserStateRegistered:
		return "REGISTERED"
	case UserStatePending:
		return "PENDING_VERIFICATION"
	case UserStateVerified:
		return "VERIFIED"
	case UserStateDeleted:
		return "DELETED"
	default:
		return ""
	}
}
func (s *UserState) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == `""` {
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch strings.ToUpper(str) {
	case "REGISTERED":
		*s = UserStateRegistered
	case "PENDING_VERIFICATION":
		*s = UserStatePending
	case "VERIFIED":
		*s = UserStateVerified
	case "DELETED":
		*s = UserStateDeleted
	default:
		return errors.New("invalid state")
	}

	return nil
}

func (s UserState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s UserState) Value() (driver.Value, error) {
	return s.String(), nil
}
func (s *UserState) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return errors.New("invalid type")
	}

	switch strings.ToUpper(str) {
	case "REGISTERED":
		*s = UserStateRegistered
	case "PENDING_VERIFICATION":
		*s = UserStatePending
	case "VERIFIED":
		*s = UserStateVerified
	case "DELETED":
		*s = UserStateDeleted
	default:
		return errors.New("invalid state")
	}

	return nil
}
