Presto for Go
=============

This is a tiny golang client for Facebook's Presto SQL Tool.

Getting Started Example:

```bash
go get github.com/tmickleydoyle/go-presto

go-presto -in query_one.sql -out query_one.json -json true

go-presto -in query_one.sql -out query_one.csv -json false

go-presto -in query_one.sql -out query_one.csv

go-presto -in query_one.sql

go-presto -in "SELECT * FROM table" -out query_one.csv

# Renders markdown table
go-presto -in "SELECT * FROM table"
```

Database connection objects should be included in the ~/.bash_profile:

```bash
export PRESTO_USERNAME='username'
export PRESTO_HOST='presto.coordinator.net'
export PRESTO_PORT=8080
```

Help:

```bash
go-presto --help

Usage of go-presto:
  -in string
        input SQL file name
  -json
        indicate if the output should be json
  -out string
        output file name
```

Go can be installed with [Homebrew](https://formulae.brew.sh/formula/go).

Install go-presto:

```bash
brew uninstall tmickleydoyle/go-presto/main
```

Install Go:

```bash
brew install go
```
