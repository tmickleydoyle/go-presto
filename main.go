package main

import (
	"fmt"
	"github.com/tmickleydoyle/go-presto"
  )
  
  // Host, user, source, catalog, schema, query
  sql := "SELECT * FROM sys.node"
  query, _ := presto.NewQuery("http://presto-coordinator:8080", "", "", "", "", sql)
  
  if row, _ := query.Next(); row != nil {
	fmt.Println(row...)
  }