package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	_tableSQL = "SELECT `table_name` AS `name`, `table_comment` AS `comment` FROM information_schema.tables WHERE table_schema = ?"
	_fieldSQL = "SELECT COLUMN_NAME name, COLUMN_KEY col_key, COLUMN_COMMENT comment, DATA_TYPE data_type FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position ASC"
)

var (
	importMap = map[string]string{
		"time.Time":             "time",
		"soft_delete.DeletedAt": "gorm.io/plugin/soft_delete",
	}
)

type Table struct {
	Name    string
	Comment string

	MOD          int      `gorm:"-"`
	PriField     *Field   `gorm:"-"`
	IsSplit      bool     `gorm:"-"`
	FirstLetter  string   `gorm:"-"`
	Imports      string   `gorm:"-"`
	PackageName  string   `gorm:"-"`
	BigCamelName string   `gorm:"-"`
	Fields       []*Field `gorm:"-"`
}

type Field struct {
	Name     string
	ColKey   string
	Comment  string
	DataType string

	BigCamelName   string `gorm:"-"`
	SmallCamelName string `gorm:"-"`
	Tag            string `gorm:"-"`
}

func GetTables() []*Table {
	connection := connectMySQL()

	var tables []*Table
	if len(Cfg.Tables) > 0 {
		connection.Raw(fmt.Sprintf("%s%s", _tableSQL, " AND table_name IN (?)"), Cfg.DBName, Cfg.Tables).Scan(&tables)
	} else {
		connection.Raw(_tableSQL, Cfg.DBName).Scan(&tables)
	}

	for i, table := range tables {
		connection.Raw(_fieldSQL, Cfg.DBName, table.Name).Scan(&tables[i].Fields)
	}
	return initTables(tables)
}

func initTables(tables []*Table) []*Table {
	res := make([]*Table, 0)
	uni := make(map[string]*Table)
	for i := range tables {
		tableName, isSplit, splitNum := parseTableName(tables[i].Name)
		if table, ok := uni[tableName]; ok {
			if !isSplit {
				panic(fmt.Sprintf("table name [%s] duplicate", tableName))
			}
			table.MOD = MaxFunc(table.MOD, splitNum)
			continue
		}

		imports := make(map[string]struct{})
		if isSplit {
			imports["fmt"] = struct{}{}
		}

		var (
			longestGormTapLen int
			priField          *Field
		)
		for _, field := range tables[i].Fields {
			tagGORMLen := len(fmt.Sprintf("column:%s", field.Name))
			if field.ColKey == "PRI" {
				priField = field
				tagGORMLen += len(";primaryKey")
			}
			longestGormTapLen = MaxFunc(longestGormTapLen, tagGORMLen)
		}

		for _, field := range tables[i].Fields {
			field.Comment = strings.ReplaceAll(field.Comment, "\n", " ")
			field.BigCamelName = ToBigCamelCase(field.Name)
			field.SmallCamelName = ToSmallCamelCase(field.Name)
			field.DataType = TransformType(field.DataType, field.Name)

			if imp, ok := importMap[field.DataType]; ok {
				imports[imp] = struct{}{}
			}
			gormTag := fmt.Sprintf("column:%s", field.Name)
			if field.ColKey == "PRI" {
				gormTag += ";primaryKey"
			}

			tag := fmt.Sprintf(`gorm:"%s"`, gormTag)
			for k := 0; k < longestGormTapLen-len(gormTag)+1; k++ {
				tag += " "
			}
			tag += fmt.Sprintf(`json:"%s"`, field.Name)
			field.Tag = fmt.Sprintf("`%s`", tag)
		}

		table := &Table{
			Name:         tableName,
			Comment:      commentString(tables[i].Comment),
			MOD:          splitNum,
			PriField:     priField,
			IsSplit:      isSplit,
			FirstLetter:  tableName[0:1],
			Imports:      importString(imports),
			PackageName:  strings.ReplaceAll(strings.ToLower(tableName), "_", ""),
			BigCamelName: ToBigCamelCase(tableName),
			Fields:       tables[i].Fields,
		}
		uni[tableName] = table
		res = append(res, table)
	}
	return res
}

func parseTableName(name string) (parsedName string, isSplit bool, splitNum int) {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] != '_' {
			continue
		}
		// 最后一个 _ 当作分表符
		num, err := strconv.Atoi(name[i+1:])
		if err != nil {
			parsedName = name
			isSplit = false
			splitNum = 0
			return
		}
		parsedName = name[:i]
		isSplit = true
		splitNum = num
		return
	}
	return name, false, 0
}

func importString(imports map[string]struct{}) string {
	if len(imports) == 0 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString("import (\n")
	for s := range imports {
		builder.WriteString(`"`)
		builder.WriteString(s)
		builder.WriteString(`"`)
		builder.WriteString("\n")
	}
	builder.WriteString(")")
	return builder.String()
}

func commentString(comment string) string {
	comment = strings.Trim(comment, " ")
	if comment == "" {
		return "."
	}
	return comment
}
