package sqlc

import (
	"bytes"
	"database/sql"
	"io"
	"fmt"
	"reflect"
	"time"
)

type PredicateType int
type JoinType int
type Dialect int

const (
	EqPredicate PredicateType = iota
	GtPredicate
	GePredicate
	LtPredicate
	LePredicate
	InPredicate
)

const (
	Join JoinType = iota
	LeftOuterJoin
	NotJoined
)

const (
	Sqlite Dialect = iota
	MySQL
	Postgres
)

type Aliasable interface {
	Alias() string
	MaybeAlias() string
}

type TableLike interface {
	Selectable
	Name() string
	As(string) Selectable
	Queryable
}

type FieldFunction struct {
	Child *FieldFunction
	Name  string
	Expr  string
	Args  []interface{}
}

type Field interface {
	Aliasable
	Functional
	Name() string
	As(string) Field
	Function() FieldFunction
}

type TableField interface {
	Field
	Parent() Selectable
}

type FieldBinding struct {
	Field Field
	Value interface{}
}

type TableFieldBinding struct {
	Field TableField
	Value interface{}
}

type Condition struct {
	Binding   FieldBinding
	Predicate PredicateType
}

type SelectFromStep interface {
	From(Selectable) SelectWhereStep
}

type SelectJoinStep interface {
	Join(Selectable) SelectOnStep
	LeftOuterJoin(Selectable) SelectOnStep
}

type SelectOnStep interface {
	On(...JoinCondition) SelectWhereStep
	Query
}

type SelectWhereStep interface {
	Query
	SelectGroupByStep
	SelectJoinStep
	Where(conditions ...Condition) Query
}

type SelectGroupByStep interface {
	GroupBy(...Field) SelectHavingStep
}

type SelectHavingStep interface {
	SelectOrderByStep
	Query
}

type SelectOrderByStep interface {
	OrderBy(...Field) SelectLimitStep
}

type SelectLimitStep interface {
	Query
}

type InsertResultStep interface {
	Renderable
	Fetch() (*sql.Row, error)
}

type InsertSetMoreStep interface {
	Executable
	InsertSetStep
	Returning(TableField) InsertResultStep
}

type UpdateSetMoreStep interface {
	Executable
	UpdateSetStep
	Where(conditions ...Condition) Executable
}

type Renderable interface {
	Render(Dialect, io.Writer) []interface{}
	String(Dialect) string
	DB() *DB
}

type Queryable interface {
	Fields() []Field
}

type Query interface {
	Renderable
	Selectable
	Query() (*sql.Rows, error)
	QueryRow() (*sql.Row, error)
}

type Executable interface {
	Renderable
	Exec() (sql.Result, error)
}

type Selectable interface {
	Aliasable
	Reflectable
	IsSelectable()
}

type JoinCondition struct {
	Lhs, Rhs  TableField
	Predicate PredicateType
}

type join struct {
	target   Selectable
	joinType JoinType
	conds    []JoinCondition
}

type update struct {
	db        *DB
	table     TableLike
	bindings  []TableFieldBinding
	predicate []Condition
}

type DB struct {
	*sql.DB
	dialect Dialect
}

func (dialect Dialect) string() string {
	switch dialect {
		case Sqlite: return "sqlite3"
		case MySQL: return "mysql"
		case Postgres: return "postgres"
	}
	return ""
}

func Open(dialect Dialect, dataSourceName string) (*DB, error) {
	db, err := sql.Open(dialect.string(), dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db, dialect}, nil
}

func (db *DB) Update(t TableLike) UpdateSetStep {
	return &update{db: db, table: t}
}

func (u *update) DB() *DB {
	return u.db
}

func (u *update) Where(c ...Condition) Executable {
	u.predicate = c
	return u
}

func (u *update) _set(f TableField, v interface{}) UpdateSetMoreStep {
	binding := TableFieldBinding{Field: f, Value: v}
	u.bindings = append(u.bindings, binding)
	return u
}

func (u *update) Exec() (sql.Result, error) {
	return exec(u)
}

func exec(r Renderable) (sql.Result, error) {
	var buf bytes.Buffer
	args := r.Render(r.DB().dialect, &buf)
	return r.DB().Exec(buf.String(), args...)
}

func (u *update) Set(f TableField, v interface{}) UpdateSetMoreStep {
	switch f := f.(type) {
		case StringField:
			_v, ok := v.(string)
			if ok {
				return u.SetString(f, _v)
			}
			panic(fmt.Sprintf("%s expects a string, got %s.", f.Name(), reflect.TypeOf(v)))

		case IntField:
			_v, ok := v.(int)
			if ok {
				return u.SetInt(f, _v)
			}
			panic(fmt.Sprintf("%s expects an int, got %s.", f.Name(), reflect.TypeOf(v)))

		case Int64Field:
			_v, ok := v.(int64)
			if ok {
				return u.SetInt64(f, _v)
			}
			panic(fmt.Sprintf("%s expects an int64, got %s.", f.Name(), reflect.TypeOf(v)))

		case TimeField:
			_v, ok := v.(time.Time)
			if ok {
	 			u.SetTime(f, _v)
	 		}
			panic(fmt.Sprintf("%s expects a time, got %s.", f.Name(), reflect.TypeOf(v)))	 		
	}
	panic("unreachable")
}

