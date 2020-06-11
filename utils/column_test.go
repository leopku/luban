package utils

import (
  "database/sql"
  // "fmt"
  "testing"

  "github.com/DATA-DOG/go-sqlmock"
  "github.com/jimsmart/schema"
  "github.com/stretchr/testify/assert"
)

var db *sql.DB
var mock sqlmock.Sqlmock

func SetUp(t *testing.T) {
  var err error
  db, mock, err = sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

}

func TestColumnGoName(t *testing.T) {
  SetUp(t)
  tables, err := schema.TableNames(db)
  assert.Nil(t, err)
  assert.NotNil(t, tables)
}
