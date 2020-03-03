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
	"strings"
	"time"

	"github.com/joho/sqltocsv"
	"github.com/jedib0t/go-pretty/table"
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

	// if outFilename != "" {
	if jsonOutput == true {
		b, err := queryToJson(rows)
		if err != nil {
			log.Fatalln(err)
		}

		jData := string(b)
		
		if outFilename != "" {
			file, _ := json.MarshalIndent(jData, "", " ")
			_ = ioutil.WriteFile(outFilename, file, 0644)
		} else {
			fmt.Println(jData)
		}
	} else {
		csvConverter := sqltocsv.New(rows)
		if outFilename != "" {
			csvConverter.WriteFile(outFilename)
		} else {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			columnNames, _ := rows.Columns()
			count := len(columnNames)
			values := make([]interface{}, count)
			valuePtrs := make([]interface{}, count)
			colNames := make([]interface{}, count)

			for i, col := range columnNames {
				colNames[i] = fmt.Sprintf("%v", col)
				}
			
			t.AppendHeader(colNames)

			for rows.Next() {
				row := make([]interface{}, count)
		
				for i, _ := range columnNames {
					valuePtrs[i] = &values[i]
				}
		
				if err = rows.Scan(valuePtrs...); err != nil {
					return err
				}
		
				for i, _ := range columnNames {
					var value interface{}
					rawValue := values[i]
		
					byteArray, ok := rawValue.([]byte)
					if ok {
						value = string(byteArray)
					} else {
						value = rawValue
					}
		
					timeValue, ok := value.(time.Time)
					if ok && csvConverter.TimeFormat != "" {
						value = timeValue.Format(csvConverter.TimeFormat)
					}
		
					if value == nil {
						row[i] = ""
					} else {
						cleanValue := fmt.Sprintf("%v", value)
						row[i] = strings.Replace(cleanValue, " ", "_", -1)
					}

					
				}
				t.AppendRow(row)
				fmt.Println(row)
			}
			t.Render()
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
