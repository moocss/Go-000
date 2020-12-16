package repository

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var mu sync.Mutex

//type SqlProvider struct {
//	DB *sql.DB
//}

// NewDB... 连接数据库
func NewDB(dsn string) (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	// 模拟数据库实例
	db, err := sql.Open("mysql", dsn)

	if err != nil || db == nil {
		return nil, err
	}

	return db, nil
}
