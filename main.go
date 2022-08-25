package main

const _prefix = "gen_"

func main() {
	initConfig()
	tables := GetTables()
	generateModel(tables)
}
