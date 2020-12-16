package repository

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/moocss/example/pkg/conf"
)

var mu sync.Mutex

type DB struct {
	DB *sql.DB
}

// NewDB... 连接数据库
func NewDB(config *conf.Conf) *DB {
	mu.Lock()
	defer mu.Unlock()

	// 模拟数据库实例
	db, err := sql.Open("mysql", config.Get("DB_DEFAULT_DSN"))

	if err != nil || db == nil {
		return nil
	}

	return &DB{DB: db}
}
