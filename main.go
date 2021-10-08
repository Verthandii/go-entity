package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

const _modelPath = "./model/"

func main() {
	initConfig()

	tables := GetTables()
	for _, table := range tables {
		modelStr, err := genGo(table)
		if err != nil {
			fmt.Println("err:", err.Error())
			return
		}

		formated, err := format.Source([]byte(modelStr))
		if err != nil {
			fmt.Println("err:", err.Error())
			return
		}

		_, err = os.Open(_modelPath)
		if err != nil {
			err = os.MkdirAll(_modelPath, 0644)
			if err != nil {
				fmt.Println("err:", err.Error())
				return
			}
		}

		err = ioutil.WriteFile(_modelPath+table.Name+".go", formated, 0644)
		if err != nil {
			fmt.Println("err:", err.Error())
			return
		}
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
