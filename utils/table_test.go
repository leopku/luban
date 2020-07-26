package utils

import (
  // "bytes"
  // "database/sql"
  "fmt"

  // "io/ioutil"
  "testing"

  // "github.com/jimsmart/schema"
  // "github.com/spf13/viper"
  "github.com/stretchr/testify/assert"
)

/*var (
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
output = "generated"
  `)
  reader := bytes.NewReader(bin)
  t.Log("reader len", reader.Len())
  vip = BuildConfig(reader, "toml")
  // db = BuildDB(reader, "toml")
  t.Log("vip", vip)
  // db = NewDB(vip)
  t.Log("db", db)
}*/

func TestGetModulePath(t *testing.T) {
  SetUp(t)
  base := vip.GetString("generation.output")
  tables, err := GetAllTableMeta(db, vip)
  if err != nil {
    t.Fatalf("an error '%s' while GetAllTableMeta", err)
  }
  t0 := tables[0]
  want := fmt.Sprintf("%s/module/%s", base, t0.GetModuleName())
  got := t0.GetModulePath(base)

  assert.Equal(t, want, got)
}

func TestGetFileName(t *testing.T) {
  SetUp(t)
  base := vip.GetString("generation.output")
  tables, err := GetAllTableMeta(db, vip)
  if err != nil {
    t.Fatalf("an error '%s' while GetAllTableMeta", err)
  }
  t0 := tables[0]
  wantModelFile := fmt.Sprintf("%s/%s_model.go", t0.GetModulePath(base), t0.GetGoFileName())
  gotModelFile := t0.GetFileName(t0.GetModulePath(base), "model")

  assert.Equal(t, wantModelFile, gotModelFile)

  wantRepoFile := fmt.Sprintf("%s/%s_repo.go", t0.GetModulePath(base), t0.GetGoFileName())
  gotRepoFile := t0.GetFileName(t0.GetModulePath(base), "repo")

  assert.Equal(t, wantRepoFile, gotRepoFile)

}
