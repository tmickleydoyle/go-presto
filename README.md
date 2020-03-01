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
```

Install with brew:

```bash
brew install tmickleydoyle/go-presto/main
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

Install Go:

```bash
brew install go
```
