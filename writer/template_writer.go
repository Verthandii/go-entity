package writer

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"text/template"
)

type TemplateWriter struct {
	t *template.Template

	c *Config
}

func (tw *TemplateWriter) Write(data interface{}, fileName string) error {
	var buf bytes.Buffer
	if err := tw.t.Execute(&buf, data); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if tw.c.Terminal {
		fmt.Println(string(formatted))
		return nil
	}

	_, err = os.Open(tw.c.Path)
	if err != nil {
		err = os.MkdirAll(tw.c.Path, 0600)
		if err != nil {
			return err
		}
	}

	filepath := fmt.Sprintf("%s/%s", tw.c.Path, fileName)
	if err = os.WriteFile(filepath, formatted, 0600); err != nil {
		return err
	}

	log.Println(fmt.Sprintf("%s complete!", filepath))

	return nil
}

func NewTemplateWriter(t *template.Template, c *Config) *TemplateWriter {
	if c == nil {
		c = &Config{
			Terminal: false,
			Path:     "./pkg/",
		}
	}

	c.Path = path.Clean(c.Path)

	return &TemplateWriter{
		t: t,
		c: c,
	}
}
