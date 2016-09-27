package sqlc

import (
	"bytes"
	"database/sql"
	"time"
	"fmt"
	"reflect"
)

type insert struct {
	db 			  *DB
	table     TableLike
	bindings  []TableFieldBinding
	returning TableField
}

func (i *insert) DB() *DB {
	return i.db
}

func (db *DB) InsertInto(t TableLike) InsertSetStep {
	return &insert{db: db, table: t}
}

func (i *insert) Exec() (sql.Result, error) {
	return exec(i)
}

func (i *insert) Returning(f TableField) InsertResultStep {
	i.returning = f
	return i
}

func (i *insert) Fetch() (*sql.Row, error) {
	var buf bytes.Buffer
	args := i.Render(i.DB().dialect, &buf)
	return i.DB().QueryRow(buf.String(), args...), nil
}

func (i *insert) _set(f TableField, v interface{}) InsertSetMoreStep {
	binding := TableFieldBinding{Field: f, Value: v}
	i.bindings = append(i.bindings, binding)
	return i
}

func (i *insert) Set(f TableField, v interface{}) InsertSetMoreStep {
	switch f := f.(type) {
		case StringField:
			_v, ok := v.(string)
			if ok {
				return i.SetString(f, _v)
			}
			panic(fmt.Sprintf("%s expects a string, got %s.", f.Name(), reflect.TypeOf(v)))

		case IntField:
			_v, ok := v.(int)
			if ok {
				return i.SetInt(f, _v)
			}
			panic(fmt.Sprintf("%s expects an int, got %s.", f.Name(), reflect.TypeOf(v)))

		case Int64Field:
			_v, ok := v.(int64)
			if ok {
				return i.SetInt64(f, _v)
			}
			panic(fmt.Sprintf("%s expects an int64, got %s.", f.Name(), reflect.TypeOf(v)))

		case TimeField:
			_v, ok := v.(time.Time)
			if ok {
	 			i.SetTime(f, _v)
	 		}
			panic(fmt.Sprintf("%s expects a time, got %s.", f.Name(), reflect.TypeOf(v)))	 		
	}
	panic("unreachable")
}

