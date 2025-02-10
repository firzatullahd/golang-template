package service

import (
	"context"
	"errors"
	"testing"

	"github.com/firzatullahd/golang-template/config"
	"github.com/firzatullahd/golang-template/internal/user/model"
	customerror "github.com/firzatullahd/golang-template/internal/user/model/error"
	mockrepo "github.com/firzatullahd/golang-template/internal/user/service/mock"
	timeutils "github.com/firzatullahd/golang-template/utils/time"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockrepo.NewMockIrepository(ctrl)
	emailAdapter := mockrepo.NewMockIEmailClient(ctrl)

	type args struct {
		in model.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *model.RegisterResponse
		wantErr error
		domock  func()
	}{
		{
			name: "error",
			args: args{
				in: model.RegisterRequest{
					Username: "test@gmail.com",
					Name:     "nama panjang",
					Password: "alalala",
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
			domock: func() {
				repo.EXPECT().WithTransaction().Return(nil, errors.New("db error")).Times(1)
			},
		},
		{
			name: "error validation email",
			args: args{
				in: model.RegisterRequest{
					Username: "test",
				},
			},
			want:    nil,
			wantErr: customerror.ErrValidationUsername,
			domock:  func() {},
		},
		{
			name: "error validation name",
			args: args{
				in: model.RegisterRequest{
					Username: "test@gmail.com",
					Name:     "a",
				},
			},
			want:    nil,
			wantErr: customerror.ErrValidationName,
			domock:  func() {},
		},
		{
			name: "error validation password",
			args: args{
				in: model.RegisterRequest{
					Username: "test@gmail.com",
					Name:     "nama panjang",
					Password: "a",
				},
			},
			want:    nil,
			wantErr: customerror.ErrValidationPassword,
			domock:  func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				conf:        &config.Config{},
				repo:        repo,
				emailClient: emailAdapter,
				time:        timeutils.Time{},
			}
			tt.domock()
			got, err := s.Register(context.Background(), tt.args.in)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
