package utils

import (
	"database/sql"
	"fmt"

	// "github.com/CloudyKit/jet"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jimsmart/schema"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	// "github.com/src-d/go-mysql-server"

	"github.com/src-d/go-mysql-server/memory"
	fakesql "github.com/src-d/go-mysql-server/sql"
)

func NewDB(vip *viper.Viper) *sql.DB {
	connStr := ""

	adapter := vip.GetString("database.adapter")
	log.Trace().Interface("vip", vip).Msg("")
	log.Trace().Str("adapter", adapter).Msg("")

	switch adapter {
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True", viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.host"), viper.GetInt("database.port"), viper.GetString("database.database"), viper.GetString("database.encoding"))
	case "postgres":
		connStr = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", viper.GetString("database.host"), viper.GetInt("database.port"), viper.GetString("database.database"), viper.GetString("database.username"), viper.GetString("database.password"))
	}
	log.Trace().Str("connection string", connStr).Msg("")

	db, err := sql.Open(adapter, connStr)
	if err != nil {
		return nil
	}
	return db
}

func GetAllTableMeta(db *sql.DB, vip *viper.Viper) ([]*TableMeta, error) {
	ret := []*TableMeta{}
	tableSlice, err := schema.TableNames(db)
	if err != nil {
		return nil, err
	}
	prefix := vip.GetString("generation.prefix")
	for _, tName := range tableSlice {
		found := false
		for _, exclude := range vip.GetStringSlice("generation.exclude") {
			if tName == exclude {
				found = true
				break
			}
		}
		if !found {
			meta := &TableMeta{
				db:         db,
				Name:       tName,
				Prefix:     prefix,
				OutputPath: vip.GetString("generation.model.output"),
			}
			ret = append(ret, meta)
		}
	}
	return ret, nil
}

func CreateTestDatabase() *memory.Database {
	const (
		dbName  = "test"
		tblName = "prefix_mytable"
	)

	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tblName, fakesql.Schema{
		{Name: "id", Type: fakesql.Int32, Nullable: false, Source: tblName},
		{Name: "name", Type: fakesql.Text, Nullable: false, Source: tblName},
		{Name: "created_at", Type: fakesql.Timestamp, Nullable: false, Source: tblName},
		{Name: "updated_at", Type: fakesql.Timestamp, Nullable: true, Source: tblName},
	})

	db.AddTable(tblName, table)

	return db
}

/*func CreateInfoSchemaDatabase() *memory.Database {
  const (
    dbName  = "information_schema"
    tblName = "TABLES"
  )

  db := memory.NewDatabase(dbName)
  table := memory.NewTable(tblName, fakesql.Schema{
    {Name: "CATALOG_NAME", Type: fakesql.Text, Nullable: false, Source: tblName},
    {Name: "SCHEMA_NAME", Type: fakesql.Text, Nullable: false, Source: tblName},
    {Name: "DEFAULT_CHARACTER_SET_NAME", Type: fakesql.Text, Nullable: false, Source: tblName},
    {Name: "DEFAULT_COLLATION_NAME", Type: fakesql.Text, Nullable: true, Source: tblName},
  })

  db.AddTable(tblName, table)

  ctx := fakesql.NewEmptyContext()

  return db
}*/
