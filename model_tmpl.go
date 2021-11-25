package main

const modelTmpl = `
{{ $TableName := .BigCamelName }}
{{ $FirstLetter := .FirstLetter }}
package model

{{ .Imports }}

// {{ $TableName }} {{ .Comment }}
type {{ $TableName }} struct {
{{ range $index, $field := .Fields }}{{ $field.BigCamelName }} {{ $field.DataType }} {{ $field.Tag }} // {{ $field.Comment }}
{{ end }}}

func (*{{ $TableName }}) TableName() string {
	return "{{ .Name }}"
}

{{ $SliceStruct := printf "%sSlice" $TableName }}

type {{ $SliceStruct }} []*{{ $TableName }}

{{ range $index, $field := .Fields }}
{{ if eq $field.ColKey "PRI" }}
func ({{ $FirstLetter }} *{{ $SliceStruct }}) IdMap() map[{{ $field.DataType }}]*{{ $TableName }} {
	uni := make(map[{{ $field.DataType }}]*{{ $TableName }})
	for _, item := range *{{ $FirstLetter }} {
		uni[item.{{ $field.BigCamelName }}] = item
	}
	return uni
}
{{ else }}

func ({{ $FirstLetter }} *{{ $SliceStruct }}) GroupBy{{ $field.BigCamelName }}() map[{{ $field.DataType }}]{{ $SliceStruct }} {
	res := make(map[{{ $field.DataType }}]{{ $SliceStruct }})
	for _, item := range *{{ $FirstLetter }} {
		res[item.{{ $field.BigCamelName }}] = append(res[item.{{ $field.BigCamelName }}], item)
    }
    return res
}
{{ end }}
{{ end }}

{{ range $index, $field := .Fields }}
func ({{ $FirstLetter }} *{{ $SliceStruct }}) Pluck{{ $field.BigCamelName }}() []{{ $field.DataType }} {
    	res := make([]{{ $field.DataType }}, 0, len(*{{ $FirstLetter }}))
    	for _, item := range *{{ $FirstLetter }} {
    		res = append(res, item.{{ $field.BigCamelName }})
    	}
    	return res
}
{{ end }}

{{ range $index, $field := .Fields }}
func ({{ $FirstLetter }} *{{ $SliceStruct }}) Unique{{ $field.BigCamelName }}() []{{ $field.DataType }} {
	uni := make(map[{{ $field.DataType }}]struct{})
	res := make([]{{ $field.DataType }}, 0)
	for _, item := range *{{ $FirstLetter }} {
		uni[item.{{ $field.BigCamelName }}] = struct{}{}
	}
	for key := range uni {
		res = append(res, key)
	}
	return res
}
{{ end }}
`