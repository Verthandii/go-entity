package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"text/template"
)

func generateDao(tables []Table) {
	t := template.New("dao_tmpl")
	t, err := t.Parse(daoTmpl)
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
			WriteFile(formated, Cfg.Output+"dao/", table.Name+".go")
		}
	}
}
