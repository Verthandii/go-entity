package main

import (
	"fmt"
	"log"
	"text/template"

	"go-entity/writer"
)

func generateModel(tables []*Table) {
	t := template.New("model_tmpl")
	t, err := t.Parse(modelTmpl)
	if err != nil {
		log.Fatalln("template parse err:", err)
	}

	t2 := template.New("model_field_tmpl")
	t2, err = t2.Parse(modelFieldTmpl)
	if err != nil {
		log.Fatalln("template parse err:", err)
	}

	for _, table := range tables {
		w := writer.NewTemplateWriter(
			t,
			&writer.Config{
				Terminal: Cfg.Terminal,
				Path:     fmt.Sprintf("%s/model/", Cfg.Output),
			},
		)
		if err = w.Write(table, fmt.Sprintf("%s%s.go", _prefix, table.Name)); err != nil {
			log.Fatalln("write occurred error:", err)
		}

		w = writer.NewTemplateWriter(
			t2,
			&writer.Config{
				Terminal: Cfg.Terminal,
				Path:     fmt.Sprintf("%s/model/%s", Cfg.Output, table.PackageName),
			},
		)
		if err = w.Write(table, fmt.Sprintf("%s%s.go", _prefix, table.Name)); err != nil {
			log.Fatalln("write occurred error:", err)
		}
	}
}
