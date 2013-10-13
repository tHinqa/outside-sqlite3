// Copyright (c) 2013 Tony Wilson. All rights reserved.
// See LICENCE file for permissions and restrictions.

package sqlite3

import (
	// "syscall"
	"fmt"
	o "github.com/tHinqa/outside"
	"testing"
)

func TestInit(t *testing.T) {
	UTF8()
	t.Log("Libversion():", Libversion())
	t.Log("Sourceid():", Sourceid())
}

var done bool

func callback(_ *Void, argc int, argv *[1000]uintptr, colName *[1000]uintptr) int {
	// func callback(_ *Void, argc int, argv **Char, colName **Char) int {
	if !done {
		done = true
		for i := 0; i < argc; i++ {
			fmt.Printf("%20s", o.CStrToString(colName[i]))
		}
		fmt.Println()
	}
	for i := 0; i < argc; i++ {
		fmt.Printf("%20s", o.CStrToString(argv[i]))
	}
	fmt.Println()
	return 0
}

func TestDb(t *testing.T) {
	var db *Sqlite3
	ret := Open(":memory:", &db)
	defer Close(db)
	if ret != 0 {
		t.Errorf("open: %s\n", Errmsg(db))
	}
	exec := func(cmd string) {
		ret := Exec(db, cmd, callback, nil, nil)
		if ret != OK {
			t.Errorf("%s: %s\n", cmd, Errmsg(db))
		}
	}
	exec(`create table test(ABC int, XYZ int);
		insert into test values(1234,5678);
		insert into test values(8765,4321);
		insert into test values(1111,9999);
		select * from test`)
}
