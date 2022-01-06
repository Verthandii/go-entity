package main

func main() {
	initConfig()
	tables := GetTables()
	generateModel(tables)
	generateQuery(tables)
	generateDAO(tables)
	generateDO(tables)
}
