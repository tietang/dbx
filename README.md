# dbx

#### 介绍
高性能数据库扩展工具,目标:
1. 简单高效，最大限度的保留原生特性的基础上，使其使用起来简单，使得开发更高效。
2. 高性能, 支持orm的基础上，最大限度的减少性能损耗，适用于高性能场景的数据库查询。
3. 自动表名和字段名称映射，默认自动把驼峰命名转换成下划线命名，比如 OrderItem -> order_item
4. 支持自定义表名和字段名称映射
 

 

#### 安装教程

1. go get -u  github.com/tietang/dbx 
 

#### 使用说明

##### 字段映射：

格式：db:"field name[,uni|id][,omitempty][-]"

 - uni|unique 字段为唯一索引字段
 - id|pk 字段为主键
 - omitempty 字段更新和写入时忽略
 - \- 字段在更新，写入、查询时忽略

Example:
```go

type Order struct {
    OrderId   int64           `db:"order_id,id"`
    Username  string          `db:"username"`
    UserId    string          `db:"user_id"`
    Amount    decimal.Decimal `db:"amount"`
    Quantity  int             `db:"quantity"`
    Status    int             `db:"status"`
    OrderType int             `db:"order_type"`
    PayStatus int             `db:"pay_status"`
    CreatedAt *time.Time      `db:"created_at,omitempty"`
    UpdatedAt time.Time       `db:"updated_at,omitempty"`
}
```

##### 打开连接：

```go
    settings := dbx.Settings{
        DriverName:      "mysql",
        User:            "root",
        Password:        "",
        Host:            "192.168.232.175:3306",
        MaxOpenConns:    10,
        MaxIdleConns:    2,
        ConnMaxLifetime: time.Minute * 30,
        Options: map[string]string{
            "charset":   "utf8",
            "parseTime": "true",
        },
    }
    db, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	
    model := &Model{}
    //插入
    rs,err=db.Insert(model)
    //更新
    rs,err=db.Update(model)
    model = &Model{Id: 10}
    //查询1个
    err=db.GetOne(model)
    err=db.Get(model, "select * from model where id=?", 10)
    q := &Model{Id: 10, Name: ""}
    var models []Model
    err=db.FindExample(q, &models)
    err=db.Find(&models, "select * from model where id <?", 20)
	
```

事务：

```go
    err := db.Tx(func(runner *dbx.TxRunner) error {
        model := &Model{}
        //插入
        rs, err = runner.Insert(model)
        //更新
        rs, err = runner.Update(model)
        model = &Model{Id: 10}
        //查询1个
        err = runner.GetOne(model)
        err = runner.Get(model, "select * from model where id=?", 10)
        q := &Model{Id: 10, Name: ""}
        var models []Model
        //查询列表
        err = runner.FindExample(q, &models)
        err = runner.Find(&models, "select * from model where id <?", 20)
    })
```


模型和表名映射：

```go
db.RegisterTable(&Order{}, "t_order")

```