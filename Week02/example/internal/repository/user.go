package repository

import (
	"context"
	"database/sql"
	"github.com/moocss/example/internal/model"
	"github.com/pkg/errors"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user := new(model.User)
	err := r.DB.QueryRow("SELECT id, username, email FROM user WHERE id=?", id).Scan(&user.Id, &user.UserName, &user.Email)

	// 对错误进行包装, 尽量使用通一的错误处理(Sentinel error), 屏蔽掉不同数据库的报错差异性
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(ErrRecordNotFound, "此用户不存在")
	}

	if err != nil {
		return nil, errors.Wrap(err, "查询用户出错了")
	}

	return user, err
}