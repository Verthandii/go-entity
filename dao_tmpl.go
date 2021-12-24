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

// 增

func (*{{ $TableName }}Dao) Create(tx *gorm.DB, v interface{}) error {
	return querypath.New{{ $TableName }}(tx).Create(v)
}

// 删

func (*{{ $TableName }}Dao) DeleteById(tx *gorm.DB, id int) error {
	return querypath.New{{ $TableName }}(tx).WhIdEq(id).Delete(&model.{{ $TableName }}{})
}

// 改

func (*{{ $TableName }}Dao) UpdateById(tx *gorm.DB, id int, v interface{}) error {
	return querypath.New{{ $TableName }}(tx).WhIdEq(id).Update(v)
}

// 查

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
`
