package service

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var mu sync.Mutex

// OpenConnection... 连接数据库
func OpenConnection() (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	// 模拟数据库实例
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")

	if err != nil || db == nil {
		return nil, err
	}

	return db, nil
}
