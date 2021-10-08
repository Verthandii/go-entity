package main

func TransformType(typeStr string) string {
	switch typeStr {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "year":
		return "int"
	case "float", "double", "decimal":
		return "float64"
	case "date", "time", "datetime", "timestamp":
		return "time.Time"
	case "char", "varchar", "tinyblob", "tinytext", "blob", "text", "mediumblob", "mediumtext", "longblob", "longtext":
		return "string"
	default:
		return "int"
	}
}
