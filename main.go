package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/prestodb/presto-go-client/presto"
)

var (
	userName = os.Getenv("PRESTO_USERNAME")
	prestoHost = os.Getenv("PRESTO_HOST")
	prestoPort = os.Getenv("PRESTO_PORT")
)

func main() {	
	dsn := "http://" + userName + "@" + prestoHost + ":" + prestoPort
	db, err := sql.Open("presto", dsn)
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Minute)

	if err == nil {
		rows, err := db.QueryContext(ctx, "SELECT * FROM tables LIMIT 10")
		if err != nil {
			log.Fatal(err)
		}
	
		for rows.Next() {
			fmt.Println(rows)
		}
	}
}