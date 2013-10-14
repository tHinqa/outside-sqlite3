// Copyright (c) 2013 Tony Wilson. All rights reserved.
// See LICENCE file for permissions and restrictions.

package driver

import (
	"database/sql"
	// "fmt"
	"./sqltest"
	sqlite3 "github.com/tHinqa/outside-sqlite3"
	"testing"
)

func init() {
	sqlite3.UTF8()
}

//TODO(t):Fix TestPreparedStmt fail
func TestSuite(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	// db, err := sql.Open("sqlite3", "a.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqltest.RunTests(t, db, sqltest.SQLITE)
}
