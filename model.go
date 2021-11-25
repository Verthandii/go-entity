package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"text/template"
)

func generateModel(tables []Table) {
	t := template.New("model_tmpl")
	t, err := t.Parse(modelTmpl)
	if err != nil {
		panic(err)
	}

	for _, table := range tables {
		var modelsBuf bytes.Buffer
		err = t.Execute(&modelsBuf, table)
		if err != nil {
			log.Fatalln("t.Execute error:", err)
		}

		formated, err := format.Source(modelsBuf.Bytes())
		if err != nil {
			log.Fatalln("format source occurred error:", err)
		}

		if Cfg.Terminal {
			fmt.Println(string(formated))
		} else {
			WriteFile(formated, Cfg.Output+"model/", table.Name+".go")
		}
	}
}
