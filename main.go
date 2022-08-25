package main

const _prefix = "gen_"

func main() {
	initConfig()
	tables := GetTables()
	generateModel(tables)
	// TODO Query 重写
	// generateQuery(tables)
}
