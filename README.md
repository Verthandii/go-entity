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
  "tables": [
    "user"
  ]
}
```

注意：当 `tables` 为空数组时，表示对所有的表生成代码。

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

func (Media) TableName() string {
	return "media"
}

type MediaSlice []*Media

func (receiver *MediaSlice) IdMap() map[string]*Media {
	m := make(map[string]*Media)
	for _, item := range *receiver {
		m[item.Id] = item
	}
	return m
}

func (receiver *MediaSlice) PluckId() []string {
	res := make([]string, 0, len(*receiver))
	for _, item := range *receiver {
		res = append(res, item.Id)
	}
	return res
}

func (receiver *MediaSlice) PluckUrl() []string {
	res := make([]string, 0, len(*receiver))
	for _, item := range *receiver {
		res = append(res, item.Url)
	}
	return res
}

func (receiver *MediaSlice) PluckType() []int {
	res := make([]int, 0, len(*receiver))
	for _, item := range *receiver {
		res = append(res, item.Type)
	}
	return res
}

func (receiver *MediaSlice) UniqueId() []string {
	m := make(map[string]struct{})
	res := make([]string, 0)
	for _, item := range *receiver {
		m[item.Id] = struct{}{}
	}
	for key := range m {
		res = append(res, key)
	}
	return res
}

func (receiver *MediaSlice) UniqueUrl() []string {
	m := make(map[string]struct{})
	res := make([]string, 0)
	for _, item := range *receiver {
		m[item.Url] = struct{}{}
	}
	for key := range m {
		res = append(res, key)
	}
	return res
}

func (receiver *MediaSlice) UniqueType() []int {
	m := make(map[int]struct{})
	res := make([]int, 0)
	for _, item := range *receiver {
		m[item.Type] = struct{}{}
	}
	for key := range m {
		res = append(res, key)
	}
	return res
}
```