package model

import (
	"fmt"
	"database/sql"
	"github.com/edwingeng/wuid/mysql/wuid"
)

/**
 * @Com www.github.com/sixstaredu
 * @Author 六星教育-shineyork老师
 */

var g *wuid.WUID

func initWuid(dsn string)  {
	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		// ...
		return db, true, nil
	}

	// Setup
	g = wuid.NewWUID("default", nil)
	_ = g.LoadH28FromMysql(newDB, "wuid")

}

func WUID() string {
	return fmt.Sprintf("%#016x", g.Next())
}
