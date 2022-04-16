package db

import (
	"database/sql"
)

func NewMySQL(url string) (*sql.DB, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// if err = db.Ping(); err != nil {
	// 	for i := 0; i < 10; i++ {
	// 		db, err = sql.Open("mysql", url)
	// 		if err == nil {
	// 			return nil, err
	// 		}
	// 		if err = db.Ping(); err == nil {
	// 			break
	// 		}
	// 		time.Sleep(time.Second*15)
	// 	}
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	return db, nil
}