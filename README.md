Presto for Go
=============

This is a tiny golang client for Facebook's Presto SQL Tool.

Getting Started Example:

```bash
go get github.com/tmickleydoyle/go-presto

go-presto -filename ~/Desktop/query_one.sql -outFilename ~/Desktop/query_one.json -jsonOutput true

go-presto -filename ~/Desktop/query_one.sql -outFilename ~/Desktop/query_one.csv -jsonOutput false

go-presto -filename ~/Desktop/query_one.sql -outFilename ~/Desktop/query_one.csv
```

Go can be installed with [Homebrew](https://formulae.brew.sh/formula/go).

Install Go:

```bash
brew install go
```
