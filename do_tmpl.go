package main

const doTmpl = `
{{ $TableName := .BigCamelName }}
package do

type {{ $TableName }} struct {
    {{ $TableName }} *model.{{ $TableName }}
}

func New{{ $TableName }}(v *model.{{ $TableName }}) *{{ $TableName }} {
	return &{{ $TableName }}{
		{{ $TableName }}: v,
	}
}

type {{ $TableName }}Slice []*{{ $TableName }}

func New{{ $TableName }}Slice(vs model.{{ $TableName }}Slice) {{ $TableName }}Slice {
	res := make({{ $TableName }}Slice, 0, len(vs))
	for _, v := range vs {
		res = append(res, New{{ $TableName }}(v))
	}
	return res
}
`
