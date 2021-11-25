package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func main() {
	initConfig()
	output := Cfg.Output

	for _, table := range GetTables() {
		modelStr, err := genGo(table)
		if err != nil {
			log.Fatalln("gen go occurred error:", err)
		}

		formated, err := format.Source([]byte(modelStr))
		if err != nil {
			log.Fatalln("format source occurred error:", err)
		}

		if Cfg.Terminal {
			fmt.Println(string(formated))
		} else {
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
}

func genGo(table Table) (string, error) {
	t := template.New("model_tmpl")
	t, err := t.Parse(modelTmpl)
	if err != nil {
		panic(err)
	}
	var modelsBuf bytes.Buffer
	err = t.Execute(&modelsBuf, table)
	if err != nil {
		return "", err
	}

	return modelsBuf.String(), nil
}
