package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func main() {
	modelPath := "./entity/"

	tables := GetTablesInfo()
	for _, table := range tables {
		modelStr, err := genGo(table)
		if err != nil {
			fmt.Println("err:", err.Error())
			return
		}
		_, err = os.Open(modelPath)
		if err != nil {
			err = os.MkdirAll(modelPath, 0644)
			if err != nil {
				fmt.Println("err:", err.Error())
				return
			}
		}

		err = ioutil.WriteFile(modelPath+ToBigCamelCase(table.Name)+".go", []byte(modelStr), 0644)
		if err != nil {
			fmt.Println("err:", err.Error())
			return
		}
	}
}

func genGo(table Table) (string, error) {
	// 解析 model
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
