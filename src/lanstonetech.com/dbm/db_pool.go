package dbm

import (
	sqlx "lanstonetech.com/dbm/database"
)

var Sessions sessions

type sessions struct {
	db map[string]*sqlx.SQLConn
}

func init() {
	Sessions.db = make(map[string]*sqlx.SQLConn)
}

func initDBSession(session string) (*sqlx.SQLConn, error) {
	result, ok := Sessions.db[session]
	if ok {
		return result, nil
	}

	db := new(sqlx.Conn)
	if err := db.Init(session); err != nil {
		panic(err)
	}

	Sessions.db[session] = db

	return db, nil
}
