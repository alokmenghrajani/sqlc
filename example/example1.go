package main

import (
  "fmt"
  "os"
  _ "github.com/mattn/go-sqlite3"
  "github.com/alokmenghrajani/sqlc/sqlc"
)
/*
type Book struct {
  ID     int  `db:"id"`
  Title  string `db:"title"`
  Author int `db:"author"`
}
*/
func main() {
  os.Remove("/tmp/example1.db")
  db, err := sqlc.Open(sqlc.Sqlite, "/tmp/example1.db")
  panicOnError(err)
  defer db.Close()

  _, err = db.Exec("CREATE TABLE books (id int primary key, title varchar(255), author int)")
  panicOnError(err)

  // Sample insert
  db.InsertInto(BOOKS).Set(BOOKS.TITLE, "Defender Of Greatness").Set(BOOKS.ID, 1).Set(BOOKS.AUTHOR, 1234).Exec()
  db.InsertInto(BOOKS).Set(BOOKS.TITLE, "Destiny Of Silver").Set(BOOKS.ID, 3).Set(BOOKS.AUTHOR, 1234).Exec()

  // Sample update
  db.Update(BOOKS).Set(BOOKS.ID, 3).Where(BOOKS.ID.Eq(2)).Exec()

  // Sample query
  rows, err := db.Select().From(BOOKS).Query()
  panicOnError(err)
  defer rows.Close()
  for rows.Next() {
    var book books
    err := rows.StructScan(&book)
    panicOnError(err)
    fmt.Printf("%d: title: %s, author: %d\n", book.Id, book.Title, book.Author)
  }
}

func panicOnError(err error) {
  if err != nil {
    panic(err)
  }
}
