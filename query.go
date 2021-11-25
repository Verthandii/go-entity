package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"text/template"
)

func generateQuery(tables []Table) {
	t := template.New("query_tmpl")
	t, err := t.Parse(queryTmpl)
	if err != nil {
		panic(err)
	}

	formated, err := format.Source([]byte(queryGo))
	if err != nil {
		log.Fatalln("format source occurred error:", err)
	}
	if Cfg.Terminal {
		fmt.Println(string(formated))
	} else {
		WriteFile(formated, Cfg.Output+"query/", "query.go")
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
			WriteFile(formated, Cfg.Output+"query/", table.Name+".go")
		}
	}
}
