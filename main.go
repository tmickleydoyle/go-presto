package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/prestodb/presto-go-client/presto"
)

var (
	userName   = os.Getenv("PRESTO_USERNAME")
	prestoHost = os.Getenv("PRESTO_HOST")
	prestoPort = os.Getenv("PRESTO_PORT")

	cancelTime = 20 * time.Minute
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("set a sql read and write path.")
		fmt.Println("Example: go-presto /Desktop/query.sql /Desktop/output.csv")
		return
	}

	filename := os.Args[1]
	outFilename := os.Args[2]

	filerc, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer filerc.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(filerc)
	f := buf.String()

	dsn := "http://" + userName + "@" + prestoHost + ":" + prestoPort
	db, err := sql.Open("presto", dsn)

	b, err := queryToJson(db, f)
	if err != nil {
		log.Fatalln(err)
	}
	// os.Stdout.Write(b)
	jString := string(b)

	df := dataframe.ReadJSON(strings.NewReader(jString))

	fmt.Println(df)

	wf, err := os.Create(outFilename)
	if err != nil {
		log.Fatal(err)
	}

	df.WriteCSV(wf)
}

func queryToJson(db *sql.DB, query string) ([]byte, error) {
	// an array of JSON objects
	// the map key is the field name
	var objects []map[string]interface{}

	ctx, cancel := context.WithTimeout(context.Background(), cancelTime)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cancel()

	for rows.Next() {
		// figure out what columns were returned
		// the column names will be the JSON object field keys
		columns, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}

		// Scan needs an array of pointers to the values it is setting
		// This creates the object and sets the values correctly
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	// indent because I want to read the output
	return json.MarshalIndent(objects, "", "\t")
}
