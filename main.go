package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func main() {
	initConfig()
	output := Cfg.Output

	tables := GetTables()
	for _, table := range tables {
		modelStr, err := genGo(table)
		if err != nil {
			log.Fatalln("gen go occurred error:", err)
		}

		formated, err := format.Source([]byte(modelStr))
		if err != nil {
			log.Fatalln("format source occurred error:", err)
		}

		_, err = os.Open(output)
		if err != nil {
			err = os.MkdirAll(output, 0644)
			if err != nil {
				log.Fatalln("os mkdir occurred error:", err)
			}
		}

		filename := output + table.Name + ".go"
		err = ioutil.WriteFile(filename, formated, 0644)
		if err != nil {
			log.Fatalln("write file occurred error:", err)
		}
		log.Println("complete file", filename)
	}
}

func genGo(table Table) (string, error) {
	modelFiles, err := template.ParseFiles("./model_tmp.tmpl")
	if err != nil {
		return "", err
	}
	var modelsBuf bytes.Buffer
	err = modelFiles.Execute(&modelsBuf, table)
	if err != nil {
		return "", err
	}

	return modelsBuf.String(), nil
}
