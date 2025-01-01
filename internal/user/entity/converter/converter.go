package converter

import (
	"github.com/firzatullahd/golang-template/internal/user/entity"
	"github.com/firzatullahd/golang-template/internal/user/model"
)

func RegisterRequestToEntity(in model.RegisterRequest, hashedPassword []byte) entity.User {
	return entity.User{
		Username: in.Username,
		Password: string(hashedPassword),
		Name:     in.Name,
	}
}
