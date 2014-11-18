package sql

import (
	"database/sql"
	"newsreader/crawler/data"
)

const (
	SQL_INSERT data.DataOperationType = iota
	SQL_UPDATE
	SQL_SELECT
)

type SqlReadWriter struct {
	InsertSqlFormat string
	UpdateSqlFormat string
	SelectSqlFormat string
	db              *sql.DB
}

func NewSqlReadWriter(db *sql.DB) *SqlReadWriter {
	return &SqlReadWriter{
		db: db,
	}
}

func (s *SqlReadWriter) Read(type_ data.DataOperationType, args ...interface{}) (interface{}, error) {
	return read(s.db, s.SelectSqlFormat, args...)
}

func (s *SqlReadWriter) Write(type_ data.DataOperationType, args ...interface{}) error {
	switch type_ {
	case SQL_INSERT:
		return write(s.db, s.InsertSqlFormat, args...)
	case SQL_UPDATE:
		return write(s.db, s.UpdateSqlFormat, args...)
	default:
		return nil
	}
}

func read(db *sql.DB, format string, args ...interface{}) (result interface{}, err error) {
	result, err = db.Query(format, args...)
	return
}

func write(db *sql.DB, format string, args ...interface{}) (err error) {
	_, err = db.Exec(format, args...)
	return
}
