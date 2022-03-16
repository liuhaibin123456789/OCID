package tool

import "database/sql"

var DB *sql.DB

func InitMysql() error {
	sql1 := "root:123456@tcp(127.0.0.1:3306)/oidc"
	db, err := sql.Open("mysql", sql1)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	DB = db
	return nil
}
