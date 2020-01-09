package main

import (
	"fmt"
	"database/sql"

	_ "github.com/prestodb/presto-go-client/presto"
)

func main() {
	dsn := "http://user@localhost:8080?catalog=default&schema=test"
	db, err := sql.Open("presto", dsn)

	db.Query("SELECT * FROM foobar WHERE id=?", 1, sql.Named("X-Presto-User", string("Alice")))

	fmt.Println(db, err)
}