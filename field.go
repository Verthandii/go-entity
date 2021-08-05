package main

type Field struct {
	Name         string
	BigCamelName string
	Comment      string
	DataType     string
	ColKey       string

	BigCamelSpaces []string
	TagGormSpaces  []string
	TagJsonSpaces  []string
	TypeSpaces     []string
}

func TransformType(typeStr string) string {
	switch typeStr {
	case "int":
		return "int64"
	case "varchar":
		return "string"
	case "timestamp":
		return "time.Time"
	default:
		return "int64"
	}
}
