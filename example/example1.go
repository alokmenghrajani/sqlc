package main

import (
  "fmt"
  "os"
  _ "github.com/mattn/go-sqlite3"
//  "database/sql"
  "github.com/alokmenghrajani/sqlc/sqlc"
)

// FOR THE README: https://github.com/square/squalor
// https://github.com/Masterminds/squirrel
// https://github.com/dropbox/godropbox/tree/master/database/sqlbuilder

func main() {
  os.Remove("./example1.db")
  db, err := sqlc.Open(sqlc.Sqlite, "./example1.db")
  panicOnError(err)
  defer db.Close()  

  _, err = db.Exec("CREATE TABLE authors (id int, name varchar(255))")

  // Sample insert
  db.InsertInto(AUTHORS).Set(AUTHORS.NAME, "toto").Set(AUTHORS.ID, 10).Exec()
  db.InsertInto(AUTHORS).Set(AUTHORS.NAME, "xyz").Set(AUTHORS.ID, 20).Exec()

  // Sample update
  db.Update(AUTHORS).Set(AUTHORS.NAME, "haha").Where(AUTHORS.ID.Eq(10)).Exec()

  // Sample query
  rows, err := db.Select(AUTHORS.NAME, AUTHORS.ID).From(AUTHORS).Query()
  panicOnError(err)
  defer rows.Close()
  for rows.Next() {
    var name string
    var id int
    err := rows.Scan(&name, &id)
    panicOnError(err)
    fmt.Printf("here: %s, %d\n", name, id)
  }
  fmt.Printf("done")
}

func panicOnError(err error) {
  if err != nil {
    panic(err)
  }
}