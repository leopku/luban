package utils

import (
  "bytes"
  "database/sql"

  // "io/ioutil"
  "testing"

  // "github.com/jimsmart/schema"
  "github.com/spf13/viper"
  "github.com/stretchr/testify/assert"
)

var (
  db  *sql.DB
  vip *viper.Viper
)

func SetUp(t *testing.T) {
  var err error

  db, err = sql.Open("mysql", "user:pass@tcp(localhost:3366)/test?charset=utf8mb4&parseTime=True")
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  // defer db.Close()
  InitConfig(
  // WithVVVV(true),
  )

  bin := []byte(`
[database]
adapter = "mysql"
host = "localhost"
port = 3366
database = "test"
username = "user"
password = "pass"

[generation]
exclude = ["prefix_ignore"]
prefix = "prefix_"
  `)
  reader := bytes.NewReader(bin)
  t.Log("reader len", reader.Len())
  vip = BuildConfig(reader, "toml")
  // db = BuildDB(reader, "toml")
  t.Log("vip", vip)
  // db = NewDB(vip)
  t.Log("db", db)
}

func TestTableCount(t *testing.T) {
  SetUp(t)
  want := 1
  tables, err := GetAllTableMeta(db, vip)

  if err != nil {
    t.Fatalf("an error '%s' while GetAllTableMeta", err)
  }

  got := len(tables)
  assert.Equal(t, want, got)
}

func TestTableName(t *testing.T) {
  SetUp(t)
  want := "prefix_mytable"
  tables, err := GetAllTableMeta(db, vip)
  if err != nil {
    t.Fatalf("an error '%s' while GetAllTableMeta", err)
  }
  t0 := tables[0]
  got := t0.Name
  assert.Equal(t, want, got)
}

func TestTableGetNameWithoutPrefix(t *testing.T) {
  SetUp(t)
  want := "mytable"
  tables, err := GetAllTableMeta(db, vip)
  if err != nil {
    t.Fatalf("an error '%s' while GetAllTableMeta", err)
  }
  t0 := tables[0]
  got := t0.GetNameWithoutPrefix()

  assert.Equal(t, want, got)
}

func TestTableGetModelName(t *testing.T) {
  SetUp(t)
  want := "Mytable"
  tables, err := GetAllTableMeta(db, vip)
  if err != nil {
    t.Fatalf("an error '%s' while GetAllTableMeta", err)
  }
  t0 := tables[0]
  got := t0.GetModelName()

  assert.Equal(t, want, got)
}

func TestTableGetGoFileName(t *testing.T) {
  SetUp(t)
  want := "mytable"
  tables, err := GetAllTableMeta(db, vip)
  if err != nil {
    t.Fatalf("an error '%s' while GetAllTableMeta", err)
  }
  t0 := tables[0]
  got := t0.GetGoFileName()

  assert.Equal(t, want, got)
}
