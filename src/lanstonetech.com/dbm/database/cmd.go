package database

import (
	"database/sql"
)

func (this *SQLConn) session(isWrite bool) (*sql.DB, error) {
	this.Lock()
	defer this.Unlock()

	return this.getconn(isWrite)
}

func (this *SQLConn) Ping() error {
	if session, err := this.session(true); err != nil || session == nil {
		return fmt.Errorf("Ping failed! session=%v,err=%v", session, err)
	} else {
		return session.Ping()
	}
}

func (this *SQLConn) isWrite(sql string) bool {

	temp := strings.ToLower(sql)
	temp = strings.TrimSpace(temp)

	if strings.HasPrefix(temp, "select") {
		return false
	}

	return true
}

func (this *SQLConn) Begin() (*sql.Tx, error) {
	db, err := this.session(true)
	if err != nil {
		return nil, err
	}

	return db.Begin()
}

///////////////////////////////////////////////////////////////////////////////
// 通用函数
///////////////////////////////////////////////////////////////////////////////
func (this *SQLConn) Prepare(query string) (*sql.Stmt, error) {

	isWrite := this.isWrite(query)

	if session, err := this.session(isWrite); err != nil || session == nil {
		return nil, fmt.Errorf("Prepare failed! query=%s,session=%v,err=%v", query, session, err)
	} else {
		return session.Prepare(this.translate(query))
	}
}

func (this *SQLConn) translate(query string) string {

	i := 1
	result := query
	for {

		idx := strings.Index(result, "?")
		if idx == -1 {
			break
		}
		result = strings.Replace(result, "?", fmt.Sprintf("$%d", i), 1)
		i++
	}

	return result
}

type QueryRowResult struct {
	err error
	row *sql.Row
}

func (r *QueryRowResult) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	return r.row.Scan(dest...)
}

func (this *SQLConn) Query(prepareSql string, args ...interface{}) (*sql.Rows, error) {

	stmtSel, err := this.Prepare(prepareSql)
	if err != nil {
		return nil, err
	}
	defer stmtSel.Close()

	rows, err2 := stmtSel.Query(args...)
	if err2 != nil {
		return nil, err2
	}
	return rows, nil
}

func (this *SQLConn) QueryRow(query string, args ...interface{}) *QueryRowResult {
	stmtSel, err := this.Prepare(query)
	if err != nil {
		return &QueryRowResult{err, nil} //保证不会崩溃
	}
	defer stmtSel.Close()
	return &QueryRowResult{nil, stmtSel.QueryRow(args...)}
}

func (this *SQLConn) Insert(prepareSql string, args ...interface{}) (int64, error) {

	var ret int64
	if err := this.QueryRow(prepareSql, args...).Scan(&ret); err != nil {
		return 0, err
	} else {
		return ret, nil
	}
}

func (this *SQLConn) Update(prepareSql string, args ...interface{}) (int64, error) {

	stmtSel, err := this.Prepare(prepareSql)
	if err != nil {
		return 0, err
	}
	defer stmtSel.Close()

	res, err := stmtSel.Exec(args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (this *SQLConn) Exec(prepareSql string, args ...interface{}) (int64, error) {

	stmtSel, err := this.Prepare(prepareSql)
	if err != nil {
		return 0, err
	}
	defer stmtSel.Close()

	res, err := stmtSel.Exec(args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
