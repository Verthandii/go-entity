package main

const querypathGo = "package querypath\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"errors\"\n\t\"time\"\n\n\t\"github.com/go-redis/redis/v8\"\n\t\"gorm.io/gorm\"\n\t\"gorm.io/gorm/clause\"\n)\n\ntype Paginate struct {\n\tPage  int   `json:\"page\"`\n\tSize  int   `json:\"size\"`\n\tTotal int64 `json:\"total\"`\n}\n\nfunc (p *Paginate) OffSet() int {\n\treturn (p.Page - 1) * p.Size\n}\n\nfunc (p *Paginate) Fix() {\n\tif p.Page <= 0 {\n\t\tp.Page = 1\n\t}\n\tif p.Size <= 0 {\n\t\tp.Size = 20\n\t}\n}\n\ntype redisWrapper struct {\n\t*redis.Client\n\tkey    string\n\texpire time.Duration\n}\n\ntype querypath struct {\n\tdb    *gorm.DB\n\tredis *redisWrapper\n\tp     *Paginate\n}\n\nfunc (q *querypath) Create(value interface{}) error {\n\tif err := q.db.Create(value).Error; err != nil {\n\t\treturn err\n\t}\n\tif q.redis != nil {\n\t\treturn q.redis.Del(context.TODO(), q.redis.key).Err()\n\t}\n\treturn nil\n}\n\nfunc (q *querypath) Delete(value interface{}) error {\n\tif err := q.db.Delete(value).Error; err != nil {\n\t\treturn err\n\t}\n\tif q.redis != nil {\n\t\treturn q.redis.Del(context.TODO(), q.redis.key).Err()\n\t}\n\treturn nil\n}\n\nfunc (q *querypath) Update(value interface{}) error {\n\tif err := q.db.Updates(value).Error; err != nil {\n\t\treturn err\n\t}\n\tif q.redis != nil {\n\t\treturn q.redis.Del(context.TODO(), q.redis.key).Err()\n\t}\n\treturn nil\n}\n\nfunc (q *querypath) Save(value interface{}) error {\n\tif err := q.db.Save(value).Error; err != nil {\n\t\treturn err\n\t}\n\tif q.redis != nil {\n\t\treturn q.redis.Del(context.TODO(), q.redis.key).Err()\n\t}\n\treturn nil\n}\n\nfunc (q *querypath) First(dest interface{}) error {\n\tif q.redis != nil {\n\t\tdata, err := q.redis.Get(context.Background(), q.redis.key).Bytes()\n\t\tif err != nil && err != redis.Nil {\n\t\t\treturn err\n\t\t}\n\t\tif err == nil && len(data) > 0 {\n\t\t\terr = json.Unmarshal(data, dest)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\treturn nil\n\t\t}\n\t}\n\n\tif err := q.db.First(dest).Error; err != nil {\n\t\treturn err\n\t}\n\n\tif q.redis != nil {\n\t\tvjson, err := json.Marshal(dest)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\tif err = q.redis.Set(context.Background(), q.redis.key, vjson, q.redis.expire).Err(); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\treturn nil\n}\n\nfunc (q *querypath) Find(dest interface{}) error {\n\tif q.p != nil {\n\t\t// 分页不走缓存\n\t\tq.db.Limit(q.p.Size).Offset(q.p.OffSet())\n\t\treturn q.db.Find(dest).Error\n\t}\n\n\tif q.redis != nil {\n\t\tdata, err := q.redis.Get(context.Background(), q.redis.key).Bytes()\n\t\tif err != nil && err != redis.Nil {\n\t\t\treturn err\n\t\t}\n\t\tif err == nil && len(data) > 0 {\n\t\t\terr = json.Unmarshal(data, dest)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\treturn nil\n\t\t}\n\t}\n\n\terr := q.db.Find(dest).Error\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tif q.redis != nil {\n\t\tvjson, err := json.Marshal(dest)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\tif err = q.redis.Set(context.Background(), q.redis.key, vjson, q.redis.expire).Err(); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\treturn nil\n}\n\nfunc (q *querypath) Count() (int64, error) {\n\tvar cnt int64\n\tif err := q.db.Count(&cnt).Error; err != nil {\n\t\treturn 0, err\n\t}\n\treturn cnt, nil\n}\n\nfunc (q *querypath) Paginate(dest interface{}) (*Paginate, error) {\n\tif q.p == nil {\n\t\treturn nil, errors.New(\"no paginate\")\n\t}\n\tq.p.Fix()\n\n\tcnt, err := q.Count()\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tq.p.Total = cnt\n\n\treturn q.p, q.Find(dest)\n}\n\nfunc (q *querypath) ForUpdate() *querypath {\n\tq.db.Clauses(clause.Locking{Strength: \"UPDATE\"})\n\treturn q\n}\n\nfunc (q *querypath) RawDB() *gorm.DB {\n\treturn q.db\n}\n"

const querypathTmpl = `
{{ $TableName := .BigCamelName }}
{{ $FirstLetter := .FirstLetter }}
package querypath

import (
	"gorm.io/gorm"
)

type {{ $TableName }} struct{
    *querypath
}

func New{{ $TableName }}(db *gorm.DB) *{{ $TableName }} {
	return &{{ $TableName }}{
		querypath: &querypath{
			db: db.Model(&model.{{ $TableName }}{}),
			p:  nil,
		},
	}
}

// ---------------
// -----自定义-----
// ---------------

// ---------------
// -----WHERE-----
// ---------------

{{ range $index, $field := .Fields }}
func ({{ $FirstLetter }} *{{ $TableName }}) Wh{{ $field.BigCamelName }}Eq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} = ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZWh{{ $field.BigCamelName }}Eq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ if eq $field.DataType "int" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "string" }} if v == "" { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "float64" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "time.Time" }} if v.IsZero() { return {{ $FirstLetter }} }{{ end }}

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} = ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) Wh{{ $field.BigCamelName }}NotEq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} != ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZWh{{ $field.BigCamelName }}NotEq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ if eq $field.DataType "int" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "string" }} if v == "" { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "float64" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "time.Time" }} if v.IsZero() { return {{ $FirstLetter }} }{{ end }}

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} != ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) Wh{{ $field.BigCamelName }}In(v []{{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} IN ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZWh{{ $field.BigCamelName }}In(v []{{ $field.DataType }}) *{{ $TableName }} {
	if len(v) == 0 { return {{ $FirstLetter }} }

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} IN ?", v)
	return {{ $FirstLetter }}
}

{{ if eq $field.DataType "string" }}
func ({{ $FirstLetter }} *{{ $TableName }}) Wh{{ $field.BigCamelName }}Like(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} LIKE ?", fmt.Sprintf("%%%s%%", v))
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZWh{{ $field.BigCamelName }}Like(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ if eq $field.DataType "int" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "string" }} if v == "" { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "float64" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "time.Time" }} if v.IsZero() { return {{ $FirstLetter }} }{{ end }}

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} LIKE ?", fmt.Sprintf("%%%s%%", v))
	return {{ $FirstLetter }}
}
{{ end }}
{{ end }}

func ({{ $FirstLetter }} *{{ $TableName }}) Where(query interface{}, args ...interface{}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where(query, args...)
	return {{ $FirstLetter }}
}

// ------------
// -----OR-----
// ------------

{{ range $index, $field := .Fields }}
func ({{ $FirstLetter }} *{{ $TableName }}) Or{{ $field.BigCamelName }}Eq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} = ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZOr{{ $field.BigCamelName }}Eq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ if eq $field.DataType "int" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "string" }} if v == "" { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "float64" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "time.Time" }} if v.IsZero() { return {{ $FirstLetter }} }{{ end }}

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} = ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) Or{{ $field.BigCamelName }}NotEq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} != ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZOr{{ $field.BigCamelName }}NotEq(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ if eq $field.DataType "int" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "string" }} if v == "" { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "float64" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "time.Time" }} if v.IsZero() { return {{ $FirstLetter }} }{{ end }}

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} != ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) Or{{ $field.BigCamelName }}In(v []{{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} IN ?", v)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZOr{{ $field.BigCamelName }}In(v []{{ $field.DataType }}) *{{ $TableName }} {
	if len(v) == 0 { return {{ $FirstLetter }} }

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} IN ?", v)
	return {{ $FirstLetter }}
}

{{ if eq $field.DataType "string" }}
func ({{ $FirstLetter }} *{{ $TableName }}) Or{{ $field.BigCamelName }}Like(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Where("{{ $field.Name }} LIKE ?", fmt.Sprintf("%%%s%%", v))
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) ZOr{{ $field.BigCamelName }}Like(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ if eq $field.DataType "int" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "string" }} if v == "" { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "float64" }} if v == 0 { return {{ $FirstLetter }} }
	{{ else if eq $field.DataType "time.Time" }} if v.IsZero() { return {{ $FirstLetter }} }{{ end }}

	{{ $FirstLetter }}.db.Where("{{ $field.Name }} LIKE ?", fmt.Sprintf("%%%s%%", v))
	return {{ $FirstLetter }}
}
{{ end }}
{{ end }}

// ------------------
// -----ORDER BY-----
// ------------------

{{ range $index, $field := .Fields }}
func ({{ $FirstLetter }} *{{ $TableName }}) OrderBy{{ $field.BigCamelName }}Desc() *{{ $TableName }} {
	{{ $FirstLetter }}.db.Order("{{ $field.Name }} DESC")
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) OrderBy{{ $field.BigCamelName }}Asc() *{{ $TableName }} {
	{{ $FirstLetter }}.db.Order("{{ $field.Name }} ASC")
	return {{ $FirstLetter }}
}
{{ end }}

func ({{ $FirstLetter }} *{{ $TableName }}) Raw(sql string, values ...interface{}) *{{ $TableName }} {
	{{ $FirstLetter }}.db.Raw(sql, values...)
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) SetPaginate(p *Paginate) *{{ $TableName }} {
	{{ $FirstLetter }}.p = p
	return {{ $FirstLetter }}
}

func ({{ $FirstLetter }} *{{ $TableName }}) WithCache(redis *redis.Client, key string, expire time.Duration) *{{ $TableName }} {
	{{ $FirstLetter }}.redis = &redisWrapper{redis, key, expire}
	return {{ $FirstLetter }}
}
`
