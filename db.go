package x

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/inflection"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB *sqlx.DB
	s  bool // TODO: will be removed
}

type Tx struct {
	Tx *sqlx.Tx
	s  bool
}

// connLifeTime must < mysql.wait_timeout
func NewDB(dsn string, connLifeTime time.Duration, s bool) (*DB, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(connLifeTime)
	db.SetMaxOpenConns(100)
	return &DB{
		DB: db,
		s:  s,
	}, nil
}

// Insert
func (db *DB) Insert(o interface{}) (int64, error) {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	var fs1 []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" {
			continue
		}
		fs = append(fs, "`"+f+"`")
		fs1 = append(fs1, ":"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if db.s {
		table = inflection.Plural(table)
	}
	s := "insert into `%s`(%s) values(%s)"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","), strings.Join(fs1, ","))
	rt, err := db.DB.NamedExec(s, o)
	if err != nil {
		return 0, err
	}
	id, err := rt.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Get
func (db *DB) Get(id int64, o interface{}) error {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if db.s {
		table = inflection.Plural(table)
	}
	s := "select * from %s where ID=? limit 1"
	s = fmt.Sprintf(s, table)
	if err := db.DB.Get(o, s, id); err != nil {
		return err
	}
	return nil
}

// Update
func (db *DB) Update(o interface{}) error {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" || f == "CreatedAt" {
			continue
		}
		fs = append(fs, "`"+f+"`=:"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if db.s {
		table = inflection.Plural(table)
	}
	s := "update `%s` set %s where ID=:ID"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","))

	_, err := db.DB.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

// Update
func (db *DB) UpdateOnly(o interface{}, only ...string) error {
	got := func(ss []string, s string) bool {
		for _, v := range ss {
			if v == s {
				return true
			}
		}
		return false
	}
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" || f == "CreatedAt" {
			continue
		}
		if !got(only, f) {
			continue
		}

		fs = append(fs, "`"+f+"`=:"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if db.s {
		table = inflection.Plural(table)
	}
	s := "update `%s` set %s where ID=:ID"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","))

	_, err := db.DB.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

// Update
func (db *DB) UpdateExcept(o interface{}, ignore ...string) error {
	got := func(ss []string, s string) bool {
		for _, v := range ss {
			if v == s {
				return true
			}
		}
		return false
	}
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" || f == "CreatedAt" {
			continue
		}
		if got(ignore, f) {
			continue
		}

		fs = append(fs, "`"+f+"`=:"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if db.s {
		table = inflection.Plural(table)
	}
	s := "update `%s` set %s where ID=:ID"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","))

	_, err := db.DB.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

// Delete
func (db *DB) Delete(o interface{}) error {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if db.s {
		table = inflection.Plural(table)
	}
	s := "delete from `%s` where ID=:ID"
	s = fmt.Sprintf(s, table)

	_, err := db.DB.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{
		Tx: tx,
		s:  db.s,
	}, nil
}

// Insert
func (tx *Tx) Insert(o interface{}) (int64, error) {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	var fs1 []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" {
			continue
		}
		fs = append(fs, "`"+f+"`")
		fs1 = append(fs1, ":"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if tx.s {
		table = inflection.Plural(table)
	}
	s := "insert into `%s`(%s) values(%s)"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","), strings.Join(fs1, ","))
	rt, err := tx.Tx.NamedExec(s, o)
	if err != nil {
		return 0, err
	}
	id, err := rt.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Get
func (tx *Tx) Get(id int64, o interface{}) error {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if tx.s {
		table = inflection.Plural(table)
	}
	s := "select * from %s where ID=? limit 1"
	s = fmt.Sprintf(s, table)
	if err := tx.Tx.Get(o, s, id); err != nil {
		return err
	}
	return nil
}

// Update
func (tx *Tx) Update(o interface{}) error {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" || f == "CreatedAt" {
			continue
		}
		fs = append(fs, "`"+f+"`=:"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if tx.s {
		table = inflection.Plural(table)
	}
	s := "update `%s` set %s where ID=:ID"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","))

	_, err := tx.Tx.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

// Update
func (tx *Tx) UpdateOnly(o interface{}, only ...string) error {
	got := func(ss []string, s string) bool {
		for _, v := range ss {
			if v == s {
				return true
			}
		}
		return false
	}
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" || f == "CreatedAt" {
			continue
		}
		if !got(only, f) {
			continue
		}

		fs = append(fs, "`"+f+"`=:"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if tx.s {
		table = inflection.Plural(table)
	}
	s := "update `%s` set %s where ID=:ID"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","))

	_, err := tx.Tx.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

// Update
func (tx *Tx) UpdateExcept(o interface{}, ignore ...string) error {
	got := func(ss []string, s string) bool {
		for _, v := range ss {
			if v == s {
				return true
			}
		}
		return false
	}
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	var fs []string
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i).Tag.Get("db")
		if f == "" || f == "-" || f == "ID" || f == "CreatedAt" {
			continue
		}
		if got(ignore, f) {
			continue
		}

		fs = append(fs, "`"+f+"`=:"+f)
	}
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if tx.s {
		table = inflection.Plural(table)
	}
	s := "update `%s` set %s where ID=:ID"
	s = fmt.Sprintf(s, table, strings.Join(fs, ","))

	_, err := tx.Tx.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

// Delete
func (tx *Tx) Delete(o interface{}) error {
	t := reflect.TypeOf(o)
	t1 := t.Elem()
	table := t1.String()[strings.LastIndex(t1.String(), ".")+1:]
	if tx.s {
		table = inflection.Plural(table)
	}
	s := "delete from `%s` where ID=:ID"
	s = fmt.Sprintf(s, table)

	_, err := tx.Tx.NamedExec(s, o)
	if err != nil {
		return err
	}
	return nil
}

func (tx *Tx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (tx *Tx) Rollback(err error) error {
	if err := tx.Tx.Rollback(); err != nil {
		return err
	}
	return err
}
