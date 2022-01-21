package main

const daoTmpl = `
{{ $TableName := .BigCamelName }}
package dao

import (
	"gorm.io/gorm"
)

type {{ $TableName }} struct{}

func New{{ $TableName }}() *{{ $TableName }} {
	return &{{ $TableName }}{}
}

func (*{{ $TableName }}) Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) error {
	return db.Transaction(fc)
}

// 增

func (*{{ $TableName }}) Create(tx *gorm.DB, v interface{}) error {
	return querypath.New{{ $TableName }}(tx).Create(v)
}

// 删

func (*{{ $TableName }}) DeleteById(tx *gorm.DB, id int) error {
	return querypath.New{{ $TableName }}(tx).WhIdEq(id).Delete(&model.{{ $TableName }}{})
}

// 改

func (*{{ $TableName }}) UpdateById(tx *gorm.DB, id int, v interface{}) error {
	return querypath.New{{ $TableName }}(tx).WhIdEq(id).Update(v)
}

// 查

func (*{{ $TableName }}) Find(q *querypath.{{ $TableName }}) (model.{{ $TableName }}Slice, error) {
	vs := make(model.{{ $TableName }}Slice, 0)
	if err := q.Find(&vs); err != nil {
		return nil, err
	}

	return vs, nil
}

func (*{{ $TableName }}) FindPagination(q *querypath.{{ $TableName }}, page, size int) (model.{{ $TableName }}Slice, *querypath.Paginate, error) {
	vs := make(model.{{ $TableName }}Slice, 0)
	p, err := q.SetPaginate(&querypath.Paginate{
		Page: page,
		Size: size,
	}).Paginate(&vs)
	if err != nil {
		return nil, nil, err
	}

	return vs, p, nil
}

func (*{{ $TableName }}) First(q *querypath.{{ $TableName }}) (*model.{{ $TableName }}, error) {
	v := &model.{{ $TableName }}{}
	if err := q.First(v); err != nil {
		return nil, err
	}

	return v, nil
}

func (*{{ $TableName }}) FirstById(tx *gorm.DB, id int) (*model.{{ $TableName }}, error) {
	v := &model.{{ $TableName }}{}
	if err := querypath.New{{ $TableName }}(tx).WhIdEq(id).First(v); err != nil {
		return nil, err
	}
	
	return v, nil
}

func (*{{ $TableName }}) Count(q *querypath.{{ $TableName }}) (int64, error) {
	return q.Count()
}
`
