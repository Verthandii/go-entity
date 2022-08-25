package main

import (
	"fmt"
	"log"
	"text/template"

	"go-entity/writer"
)

func generateQuery(tables []*Table) {
	t := template.New("query_tmpl")
	w := writer.NewTemplateWriter(
		t,
		&writer.Config{
			Terminal: Cfg.Terminal,
			Path:     fmt.Sprintf("%s/query/", Cfg.Output),
		},
	)

	t, err := t.Parse(queryGo)
	if err != nil {
		log.Fatalln("template parse err:", err)
	}
	if err = w.Write(nil, "query.go"); err != nil {
		log.Fatalln("write occurred error:", err)
	}

	t, err = t.Parse(queryTmpl)
	if err != nil {
		log.Fatalln("template parse err:", err)
	}
	for _, table := range tables {
		if err = w.Write(table, fmt.Sprintf("%s%s.go", _prefix, table.Name)); err != nil {
			log.Fatalln("write occurred error:", err)
		}
	}
}
