/*
db.go : GOWS database program
Copyright (C) 2013 Shaun Savage <savages@savages.com>

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Lesser General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of  MERCHANTABILITY or FITNESS FOR
A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License along with
this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package db

import (
    "github.com/jmoiron/sqlx"
	"log"
	"fmt"
	"database/sql"
)
import _ "github.com/jbarham/gopgsqldriver"

var dbsql *sqlx.DB



func GetDB() *sqlx.DB {
	return dbsql
}

func SetDB( db *sqlx.DB) {
	dbsql = db
}

func ConvNullStr( val sql.NullString) string {
	var ret string
	if val.Valid {
		ret = val.String
	} else {
		ret =  ""
	}
	return ret
}

func ConvNullInt( val sql.NullInt64) int64 {
	var ret int64
	if val.Valid {
		ret = val.Int64
	} else {
		ret =  0
	}
	return ret
}

func ConvNullFloat( val sql.NullFloat64) float64 {
	var ret float64
	if val.Valid {
		ret = val.Float64
	} else {
		ret =  0
	}
	return ret
}

func InitDB(dbname string, user string, passwd string) *sqlx.DB {
	db , err := sqlx.Open("postgres", "dbname="+dbname+" user="+user+" password="+passwd+" host=localhost" )
	if err != nil {
		log.Fatalf("Error opening db open: %s\n", err)
		return nil
	}

	// keep connection open
	// defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Printf("Error Ping: %s\n", err)
//		db.Close()
//		return nil
	}
	fmt.Println("DB Ping good")

	// Get a configuration
	gofwconf := &GofwConf{}
	fmt.Println(gofwconf)
	type idx struct {
		Max int32
	}
	idxmax := idx{}
	//err = db.Get(&gofwconf,"SELECT * FROM gofwconf WHERE cfgidx = (select max(cfgidx) from gofwconf)")
	err = db.Get(&idxmax,"SELECT max(cfgidx) FROM gofwconf")
	if err != nil {
		log.Printf("Error getting config: %s\n", err)
		fmt.Println(gofwconf)
//		db.Close()
//		return nil
	} else {
		fmt.Println(idxmax)
		fmt.Println(gofwconf.Name, gofwconf.Uuid)
	}
	dbsql = db
	return db
}
