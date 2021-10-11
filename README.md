# Why to use go-entity

解脱双手，一键生成 gorm 所需要的 entity。

# How to use

## Configuration

参照 `db.json.example` 中的格式编写配置文件，保存文件名为`db.json`。

```json
{
  "host": "127.0.0.1",
  "user": "root",
  "pwd": "pwd",
  "dbname": "test",
  "collation": "utf8mb4_unicode_ci",
  "tables": [
    "user"
  ],
  "output": "./model",
  "terminal": true
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
	Id   string `gorm:"id;PRIMARY_KEY" json:"id"`   //
	Url  string `gorm:"url"            json:"url"`  //
	Type int    `gorm:"type"           json:"type"` // 1. 图片
}

func (*Media) TableName() string {
	return "media"
}

type MediaSlice []*Media

func (m *MediaSlice) IdMap() map[string]*Media {
	uni := make(map[string]*Media)
	for _, item := range *m {
		uni[item.Id] = item
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

func (m *MediaSlice) PluckId() []string {
	res := make([]string, 0, len(*m))
	for _, item := range *m {
		res = append(res, item.Id)
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

func (m *MediaSlice) UniqueId() []string {
	uni := make(map[string]struct{})
	res := make([]string, 0)
	for _, item := range *m {
		uni[item.Id] = struct{}{}
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