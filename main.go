package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/joho/sqltocsv"
	_ "github.com/prestodb/presto-go-client/presto"
)

var (
	userName   = os.Getenv("PRESTO_USERNAME")
	prestoHost = os.Getenv("PRESTO_HOST")
	prestoPort = os.Getenv("PRESTO_PORT")

	cancelTime = 20 * time.Minute
)

func main() {
	var jsonOutput bool
	var filename string
	var outFilename string

	flag.BoolVar(&jsonOutput, "jsonOutput", false, "indicate if the output should be json")
	flag.StringVar(&filename, "filename", "", "input SQL file name")
	flag.StringVar(&outFilename, "outFilename", "", "output file name")

	flag.Parse()

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

	ctx, cancel := context.WithTimeout(context.Background(), cancelTime)
	defer cancel()

	rows, err := db.QueryContext(ctx, f)

	if jsonOutput == true {
		b, err := queryToJson(rows)
		if err != nil {
			log.Fatalln(err)
		}

		jData := string(b)
		file, _ := json.MarshalIndent(jData, "", " ")
		_ = ioutil.WriteFile(outFilename, file, 0644)
	} else {
		csvConverter := sqltocsv.New(rows)
		csvConverter.WriteFile(outFilename)
	}
}

func queryToJson(rows *sql.Rows) ([]byte, error) {
	var objects []map[string]interface{}

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
