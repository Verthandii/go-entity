package main

const daoTmpl = `
{{ $TableName := .BigCamelName }}
package dao

import (
	"gorm.io/gorm"
)

type {{ $TableName }}Dao struct{}

func New{{ $TableName }}Dao() *{{ $TableName }}Dao {
	return &{{ $TableName }}Dao{}
}

func (*{{ $TableName }}Dao) Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) error {
	return db.Transaction(fc)
}

func (*{{ $TableName }}Dao) List(q *querypath.{{ $TableName }}) (model.{{ $TableName }}Slice, error) {
	vs := make(model.{{ $TableName }}Slice, 0)
	if err := q.Find(&vs); err != nil {
		return nil, err
	}

	return vs, nil
}

func (*{{ $TableName }}Dao) First(q *querypath.{{ $TableName }}) (*model.{{ $TableName }}, error) {
	v := &model.{{ $TableName }}{}
	if err := q.First(v); err != nil {
		return nil, err
	}

	return v, nil
}

func (*{{ $TableName }}Dao) Delete(q *querypath.{{ $TableName }}) error {
	return q.Delete(&model.{{ $TableName }}{})
}

func (*{{ $TableName }}Dao) Create(q *querypath.{{ $TableName }}, v interface{}) error {
	return q.Create(v)
}

func (*{{ $TableName }}Dao) Update(q *querypath.{{ $TableName }}, v interface{}) error {
	return q.Update(v)
}
`
