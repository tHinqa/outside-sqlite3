// Copyright (c) 2013 Tony Wilson. All rights reserved.
// See LICENCE file for permissions and restrictions.

//Package driver implements database/sql/driver support for sqlite3.
/*
Includes some unchanged code, logic and methodology based 
on the package https://github.com/mattn/go-sqlite3 that is:
	Copyright 2012-2013 Yasuhiro Matsumoto <mattn.jp@gmail.com>
	License (MIT): http://mattn.mit-license.org/2012
*/
package driver

import (
	"database/sql"
	"database/sql/driver"
	o "github.com/tHinqa/outside"
	sqlite3 "github.com/tHinqa/outside-sqlite3"
	"io"
	"strings"
	"time"
	"unsafe"
)

func init() {
	sql.Register("sqlite3", &Driver{})
}

type Driver struct{}

type Conn struct {
	db *sqlite3.Sqlite3
}

type Tx struct {
	c *Conn
}

type Stmt struct {
	c      *Conn
	stmt   *sqlite3.Stmt
	tail   string
	closed bool
}

type Rows struct {
	s        *Stmt
	nCols    int
	cols     []string
	decltype []string
}

type errNum int

func (err errNum) Error() string {
	return sqlite3.Errstr(int(err))
}

func (c *Conn) exec(cmd string) error {
	ret := sqlite3.Exec(c.db, cmd, nil, nil, nil)
	if ret != sqlite3.OK {
		return errNum(ret)
	}
	return nil
}

func (c *Conn) Begin() (tx driver.Tx, err error) {
	err = c.exec("begin")
	if err == nil {
		tx = &Tx{c}
	}
	return
}

func (tx *Tx) Commit() (err error) {
	return tx.c.exec("commit")
}

func (c *Conn) Close() error {
	s := sqlite3.NextStmt(c.db, nil)
	for s != nil {
		sqlite3.Finalize(s)
		s = sqlite3.NextStmt(c.db, nil)
	}
	ret := sqlite3.Close(c.db)
	if ret != sqlite3.OK {
		return errNum(ret)
	}
	c.db = nil
	return nil
}

func (tx *Tx) Rollback() error {
	return tx.c.exec("rollback")
}

func (c *Conn) Prepare(query string) (s driver.Stmt, e error) {
	var stmt *sqlite3.Stmt
	var tail *sqlite3.Char
	ret := sqlite3.PrepareV2(c.db, query, -1, &stmt, &tail)
	if ret == sqlite3.OK {
		var t string
		if tail != nil && *tail != 0 {
			t = strings.TrimSpace(o.CStrToString((uintptr)(unsafe.Pointer(tail))))
		}
		s = &Stmt{c: c, stmt: stmt, tail: t}
	} else {
		e = errNum(ret)
	}
	return
}

func (d *Driver) Open(file string) (conn driver.Conn, err error) {
	var db *sqlite3.Sqlite3
	// ret := sqlite3.Open(o.VString(file), &db)
	ret := sqlite3.OpenV2(file, &db, sqlite3.OPEN_FULLMUTEX|
		sqlite3.OPEN_READWRITE|sqlite3.OPEN_CREATE, "")
	if ret != sqlite3.OK {
		err = errNum(ret)
	} else if db != nil {
		conn = &Conn{db}
		ret = sqlite3.BusyTimeout(db, 5000)
	} else {
		//TODO(T): no memory error
	}
	return
}

type Result struct {
	lastId   int64
	affected int64
}

func (res *Result) LastInsertId() (int64, error) {
	return res.lastId, nil
}

func (res *Result) RowsAffected() (int64, error) {
	return res.affected, nil
}

func (s *Stmt) Exec(args []driver.Value) (res driver.Result, err error) {
	if err = s.bind(args); err == nil {
		ret := sqlite3.Step(s.stmt)
		switch ret {
		case sqlite3.ROW, sqlite3.OK, sqlite3.DONE:
			res = &Result{
				sqlite3.LastInsertRowid(s.c.db),
				int64(sqlite3.Changes(s.c.db)),
			}
		default:
			err = errNum(ret)
		}
	}
	return
}

func (s *Stmt) bind(args []driver.Value) (err error) {
	ret := sqlite3.Reset(s.stmt)
	switch ret {
	case sqlite3.ROW, sqlite3.OK, sqlite3.DONE:
		for i, v := range args {
			n := i + 1
			switch v := v.(type) {
			case nil:
				ret = sqlite3.BindNull(s.stmt, n)
			case string:
				ret = sqlite3.BindText(s.stmt, n, o.VString(v), len(v), nil)
			case int:
				ret = sqlite3.BindInt64(s.stmt, n, int64(v))
			case int64:
				ret = sqlite3.BindInt64(s.stmt, n, v)
			case int32:
				ret = sqlite3.BindInt(s.stmt, n, int(v))
			case byte:
				ret = sqlite3.BindInt(s.stmt, n, int(v))
			case bool:
				b := 0
				if bool(v) {
					b = 1
				}
				ret = sqlite3.BindInt(s.stmt, n, b)
			case float32:
				ret = sqlite3.BindDouble(s.stmt, n, float64(v))
			case float64:
				ret = sqlite3.BindDouble(s.stmt, n, v)
			case []byte:
				var p *byte
				if len(v) > 0 {
					p = &v[0]
				}
				// TODO(t): Is there a reason for unsafe.Pointer
				// ret = sqlite3.BindBlob(s.stmt, n, unsafe.Pointer(p), len(v), nil)
				ret = sqlite3.BindBlob(s.stmt, n, p, len(v), nil)
			case time.Time:
				b := v.UTC().Format(TimestampFormats[0])
				ret = sqlite3.BindText(s.stmt, n, o.VString(b), len(b), nil)
			}
			if ret != sqlite3.OK {
				return errNum(ret)
			}
		}
	default:
		return errNum(ret)
	}
	return
}

func (s *Stmt) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	if s.c == nil || s.c.db == nil {
		// return errors.New("sqlite statement with already closed database connection")
	}
	ret := sqlite3.Finalize(s.stmt)
	if ret != sqlite3.OK {
		return errNum(ret)
	}
	return nil
}

func (s *Stmt) NumInput() int {
	return sqlite3.BindParameterCount(s.stmt)
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	if err := s.bind(args); err != nil {
		return nil, err
	}
	return &Rows{s, sqlite3.ColumnCount(s.stmt), nil, nil}, nil
}

func (rs *Rows) Close() error {
	if rs.s.closed {
		return nil
	}
	ret := sqlite3.Reset(rs.s.stmt)
	if ret != sqlite3.OK {
		return errNum(ret)
	}
	return nil
}

func (rs *Rows) Columns() []string {
	if rs.nCols != len(rs.cols) {
		rs.cols = make([]string, rs.nCols)
		for i := 0; i < rs.nCols; i++ {
			rs.cols[i] = string(sqlite3.ColumnName(rs.s.stmt, i))
		}
	}
	return rs.cols
}

func (rs *Rows) Next(dest []driver.Value) error {
	ret := sqlite3.Step(rs.s.stmt)
	if ret == sqlite3.DONE {
		return io.EOF
	}
	if ret != sqlite3.ROW {
		ret = sqlite3.Reset(rs.s.stmt)
		if ret != sqlite3.OK {
			return errNum(ret)
		}
		return nil
	}
	if rs.decltype == nil {
		rs.decltype = make([]string, rs.nCols)
		for i := 0; i < rs.nCols; i++ {
			rs.decltype[i] = strings.ToLower(
				string(sqlite3.ColumnDecltype(rs.s.stmt, i)))
		}
	}
	for i := range dest {
		switch sqlite3.ColumnType(rs.s.stmt, i) {
		case sqlite3.INTEGER:
			val := sqlite3.ColumnInt64(rs.s.stmt, i)
			switch rs.decltype[i] {
			case "timestamp", "datetime":
				dest[i] = time.Unix(val, 0)
			case "boolean":
				dest[i] = val > 0
			default:
				dest[i] = val
			}
		case sqlite3.FLOAT:
			dest[i] = float64(sqlite3.ColumnDouble(rs.s.stmt, i))
		case sqlite3.BLOB:
			p := sqlite3.ColumnBlob(rs.s.stmt, i)
			n := sqlite3.ColumnBytes(rs.s.stmt, i)
			switch dest[i].(type) {
			case sql.RawBytes:
				dest[i] = (*[1 << 30]byte)(unsafe.Pointer(p))[0:n:n]
			default:
				slice := make([]byte, n)
				copy(slice[:], (*[1 << 30]byte)(unsafe.Pointer(p))[0:n:n])
				dest[i] = slice
			}
		case sqlite3.NULL:
			dest[i] = nil
		case sqlite3.TEXT:
			s := string(sqlite3.ColumnText(rs.s.stmt, i))
			switch rs.decltype[i] {
			case "timestamp", "datetime":
				var err error
				for _, format := range TimestampFormats {
					if dest[i], err = time.Parse(format, s); err == nil {
						break
					}
				}
				if err != nil {
					dest[i] = time.Time{}
				}
			default:
				dest[i] = s
			}

		}
	}
	return nil
}

var TimestampFormats = []string{
	"2006-01-02 15:04:05.999999999",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
	"2006-01-02",
}
