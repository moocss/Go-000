package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/moocss/example/internal/model"
	"github.com/moocss/example/internal/repository"
	"github.com/moocss/example/pkg/errcode"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)

	GetRepository() repository.UserRepository
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (srv *userService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := srv.repo.GetByID(ctx, id)

	// 中间层 service 尽量不处理 dao 的 error, 直接透传到它的最上层.
	// 对是否有记录进行判断, 根据业务需求, 可进行更多处理
	if err != nil && errors.Is(err, errcode.ErrRecordNotFound) {
		// ...
		return nil, err
	}
	return user, nil
}

func (srv *userService) GetRepository() repository.UserRepository {
	return srv.repo
}
