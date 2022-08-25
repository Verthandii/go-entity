# Why to use go-entity

解脱双手，一键生成 `gorm` 所需要的 `entity`。

# How to use

## Configuration

参照 `db.json.example` 中的格式编写配置文件，保存文件名为`db.json`。

表名为 `Xxx_1` 时，会生成可分表的 `model`。(注意需要将数字**最大**的表名写上，不然会导致分表 `mod` 错误。)

字段名为 `deleted_at` 时，此字段的类型会变成 `soft_delete.DeletedAt`。（暂时只支持 `int` 类型的软删除。下次更新支持 `time.Time` 的软删除。）

```json
{
  "host": "127.0.0.1",
  "port": 3306,
  "username": "root",
  "password": "root",
  "dbname": "test",
  "tables": [
    "user"
  ],
  "output": "./pkg/",
  "terminal": false
}
```

> tables 为空数组时，表示对所有的表生成代码
>
> output 指定输出目录，只有当 terminal 为 false 时生效
>
> terminal 为 true 时输出到控制台

## Get Started

```shell script
go get
```

```shell script
go run .
```

## Example

```sql
create table media
(
    id   varchar(32) not null,
    url  varchar(256) null,
    type int         not null comment '1. 图片',
    constraint media_id_uindex
        unique (id)
) comment '媒体表';

alter table media
    add primary key (id);
```

```go
package model

// Media 媒体表
type Media struct {
	ID   string `gorm:"id;primaryKey" json:"id"`    //
	Url  string `gorm:"url"            json:"url"`  //
	Type int    `gorm:"type"           json:"type"` // 1. 图片
}

func (*Media) TableName() string {
	return "media"
}

type MediaSlice []*Media

func (m *MediaSlice) IDMap() map[string]*Media {
	uni := make(map[string]*Media)
	for _, item := range *m {
		uni[item.ID] = item
	}
	return uni
}

func (m *MediaSlice) GroupByUrl() map[string]MediaSlice {
	res := make(map[string]MediaSlice)
	for _, item := range *m {
		res[item.Url] = append(res[item.Url], item)
	}
	return res
}

func (m *MediaSlice) GroupByType() map[int]MediaSlice {
	res := make(map[int]MediaSlice)
	for _, item := range *m {
		res[item.Type] = append(res[item.Type], item)
	}
	return res
}

func (m *MediaSlice) PluckID() []string {
	res := make([]string, 0, len(*m))
	for _, item := range *m {
		res = append(res, item.ID)
	}
	return res
}

func (m *MediaSlice) PluckUrl() []string {
	res := make([]string, 0, len(*m))
	for _, item := range *m {
		res = append(res, item.Url)
	}
	return res
}

func (m *MediaSlice) PluckType() []int {
	res := make([]int, 0, len(*m))
	for _, item := range *m {
		res = append(res, item.Type)
	}
	return res
}

func (m *MediaSlice) UniqueID() []string {
	uni := make(map[string]struct{})
	res := make([]string, 0)
	for _, item := range *m {
		uni[item.ID] = struct{}{}
	}
	for key := range uni {
		res = append(res, key)
	}
	return res
}

func (m *MediaSlice) UniqueUrl() []string {
	uni := make(map[string]struct{})
	res := make([]string, 0)
	for _, item := range *m {
		uni[item.Url] = struct{}{}
	}
	for key := range uni {
		res = append(res, key)
	}
	return res
}

func (m *MediaSlice) UniqueType() []int {
	uni := make(map[int]struct{})
	res := make([]int, 0)
	for _, item := range *m {
		uni[item.Type] = struct{}{}
	}
	for key := range uni {
		res = append(res, key)
	}
	return res
}
```