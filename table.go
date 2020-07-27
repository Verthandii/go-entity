package main

var (
	TableSQL = "SELECT " +
		"table_name name " +
		"FROM " +
		"information_schema.tables " +
		"WHERE " +
		"table_schema=?"
	FieldSQL = "SELECT " +
		"COLUMN_NAME name, COLUMN_KEY col_key, COLUMN_COMMENT comment, DATA_TYPE data_type " +
		"FROM " +
		"information_schema.columns " +
		"WHERE " +
		"table_schema=? AND table_name=? " +
		"ORDER BY ordinal_position ASC"
)

type Table struct {
	Name         string
	BigCamelName string
	Fields       []Field

	LongestBigCamelColLen int
	LongestTagGORMLen     int
	LongestTagJsonLen     int
	LongestTypeLen        int
}

func GetTablesInfo() []Table {
	connection := Connect()

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
	return InitTables(tables)
}

func InitTables(tables []Table) []Table {
	for i := range tables {
		longestBigCamelColLen, longestTagGORMLen, longestTagJsonLen, longestTypeLen := 0, 0, 0, 0

		tables[i].BigCamelName = ToBigCamelCase(tables[i].Name)
		for j := range tables[i].Fields {
			tables[i].Fields[j].BigCamelName = ToBigCamelCase(tables[i].Fields[j].Name)
			tables[i].Fields[j].DataType = TransformType(tables[i].Fields[j].DataType)

			tagGORMLen := len(tables[i].Fields[j].Name)
			if tables[i].Fields[j].ColKey == "PRI" {
				tagGORMLen += len(";PRIMARY_KEY")
			}

			longestBigCamelColLen = MaxFunc(longestBigCamelColLen, len(tables[i].Fields[j].BigCamelName))
			longestTagGORMLen = MaxFunc(longestTagGORMLen, tagGORMLen)
			longestTagJsonLen = MaxFunc(longestTagJsonLen, len(tables[i].Fields[j].Name))
			longestTypeLen = MaxFunc(longestTypeLen, len(tables[i].Fields[j].DataType))
		}
		tables[i].LongestBigCamelColLen = longestBigCamelColLen
		tables[i].LongestTagGORMLen = longestTagGORMLen
		tables[i].LongestTagJsonLen = longestTagJsonLen
		tables[i].LongestTypeLen = longestTypeLen
	}

	for i := range tables {
		for j := range tables[i].Fields {
			tagGORMLen := len(tables[i].Fields[j].Name)
			if tables[i].Fields[j].ColKey == "PRI" {
				tagGORMLen += len(";PRIMARY_KEY")
			}

			tables[i].Fields[j].BigCamelSpaces = make([]string, tables[i].LongestBigCamelColLen-len(tables[i].Fields[j].BigCamelName)+1)
			tables[i].Fields[j].TagGormSpaces = make([]string, tables[i].LongestTagGORMLen-tagGORMLen+1)
			tables[i].Fields[j].TagJsonSpaces = make([]string, tables[i].LongestTagJsonLen-len(tables[i].Fields[j].Name)+1)
			tables[i].Fields[j].TypeSpaces = make([]string, tables[i].LongestTypeLen-len(tables[i].Fields[j].DataType)+1)
		}
	}
	return tables
}
