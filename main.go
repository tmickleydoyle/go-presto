package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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

	cancelTime = 1440 * time.Minute

	jsonBool bool
	in       string
	out      string
)

const (
	exitFail = 1
)

func main() {
	flag.BoolVar(&jsonBool, "json", false, "indicate if the output should be json")
	flag.StringVar(&in, "in", "", "input SQL file")
	flag.StringVar(&out, "out", "", "output file")
	flag.Parse()

	if err := run(jsonBool, in, out); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(jsonOutput bool, filename string, outFilename string) error {
	// if filename == "" {
	// 	return errors.New("no input SQL file")
	// }
	var f string

	if _, err := os.Stat(filename); err == nil {
		filerc, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer filerc.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(filerc)
		f = buf.String()
	} else {
		f = filename
	}

	dsn := "http://" + userName + "@" + prestoHost + ":" + prestoPort
	db, err := sql.Open("presto", dsn)

	if err != nil {
		return errors.New("no database connection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cancelTime)
	defer cancel()

	rows, err := db.QueryContext(ctx, f)

	db.Close()

	if err != nil {
		return errors.New(err.Error())
	}

	if outFilename != "" {
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

	return nil
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
