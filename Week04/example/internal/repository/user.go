package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/moocss/example/internal/model"
	"github.com/moocss/example/pkg/errcode"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := new(model.User)
	err := r.DB.QueryRow("SELECT id, username, email FROM user WHERE id=?", id).Scan(&user.Id, &user.UserName, &user.Email)

	// 对错误进行包装, 尽量使用通一的错误处理(Sentinel error), 屏蔽掉不同数据库的报错差异性
	// 这里不只是返回这一种错误, 在上层打日志的时候还要看到底层出的是那种错误, 所以原始的错误信息也需要保留.
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(errcode.ErrRecordNotFound, fmt.Sprintf("此用户不存在, err: %v", err))
	}

	if err != nil {
		return nil, errors.Wrap(err, "查询用户出错了")
	}

	return user, err
}