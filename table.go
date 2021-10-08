package main

import (
	"fmt"
	"strings"
)

var (
	TableSQL = "SELECT " +
		"table_name name " +
		"FROM " +
		"information_schema.tables " +
		"WHERE " +
		"table_schema = ?"
	FieldSQL = "SELECT " +
		"COLUMN_NAME name, COLUMN_KEY col_key, COLUMN_COMMENT comment, DATA_TYPE data_type " +
		"FROM " +
		"information_schema.columns " +
		"WHERE " +
		"table_schema = ? AND table_name = ? " +
		"ORDER BY ordinal_position ASC"
)

type Table struct {
	Name string

	Imports      string
	BigCamelName string
	Fields       []Field
}

type Field struct {
	Name     string
	ColKey   string
	Comment  string
	DataType string

	BigCamelName string
	Tag          string
}

func GetTables() []Table {
	connection := connectMySQL()

	var tables []Table
	if len(Cfg.Tables) > 0 {
		TableSQL += " AND table_name IN (?)"
		connection.DB.Raw(TableSQL, Cfg.DbName, Cfg.Tables).Scan(&tables)
	} else {
		connection.DB.Raw(TableSQL, Cfg.DbName).Scan(&tables)
	}

	// 查出每个表的字段信息
	for i, table := range tables {
		var columns []Field
		connection.DB.Raw(FieldSQL, Cfg.DbName, table.Name).Scan(&columns)
		tables[i].Fields = columns
	}
	return initTables(tables)
}

func initTables(tables []Table) []Table {
	for i := range tables {
		longestTagGORMLen := 0

		tables[i].BigCamelName = ToBigCamelCase(tables[i].Name)

		for j := range tables[i].Fields {
			tables[i].Fields[j].Comment = strings.ReplaceAll(tables[i].Fields[j].Comment, "\n", " ")
			tables[i].Fields[j].BigCamelName = ToBigCamelCase(tables[i].Fields[j].Name)
			tables[i].Fields[j].DataType = TransformType(tables[i].Fields[j].DataType)

			if tables[i].Fields[j].DataType == "time.Time" {
				tables[i].Imports = `import "time"`
			}

			tagGORMLen := len(tables[i].Fields[j].Name)
			if tables[i].Fields[j].ColKey == "PRI" {
				tagGORMLen += len(";PRIMARY_KEY")
			}

			longestTagGORMLen = MaxFunc(longestTagGORMLen, tagGORMLen)
		}

		for j := range tables[i].Fields {
			gormTag := tables[i].Fields[j].Name
			if tables[i].Fields[j].ColKey == "PRI" {
				gormTag += ";PRIMARY_KEY"
			}

			tag := fmt.Sprintf(`gorm:"%s"`, gormTag)
			for i := 0; i < longestTagGORMLen-len(gormTag)+1; i++ {
				tag += " "
			}
			tag += fmt.Sprintf(`json:"%s"`, tables[i].Fields[j].Name)

			tables[i].Fields[j].Tag = tag
		}
	}

	return tables
}
