package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewSqlConn(driver string) (*sql.DB, error) {
	sqlConn, err := sql.Open(driver, "test_user:test_pass@tcp(127.0.0.1:3306)/test")
	return sqlConn, err
}
