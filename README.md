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



### 性能比较


```
hood
                   Insert:   2000    18.30s      9149968 ns/op   12090 B/op    199 allocs/op
      MultiInsert 100 row:    500     Not support multi insert
                   Update:   2000    19.83s      9914022 ns/op   12081 B/op    199 allocs/op
                     Read:   4000    35.64s      8910711 ns/op    4242 B/op     55 allocs/op
      MultiRead limit 100:   2000    32.98s     16491001 ns/op  232327 B/op   8765 allocs/op
raw
                   Insert:   2000     7.36s      3681234 ns/op     552 B/op     12 allocs/op
      MultiInsert 100 row:    500    17.59s     35184280 ns/op  110864 B/op    811 allocs/op
                   Update:   2000     7.09s      3543659 ns/op     616 B/op     14 allocs/op
                     Read:   4000    13.27s      3317756 ns/op    1432 B/op     37 allocs/op
      MultiRead limit 100:   2000    20.48s     10241508 ns/op   34704 B/op   1320 allocs/op
qbs
                   Insert:   2000     no primary key field
      MultiInsert 100 row:    500     Not support multi insert
                   Update:   2000     no primary key field
                     Read:   4000     no primary key field
      MultiRead limit 100:   2000     no primary key field
gorp
Error 1054: Unknown column 'id,pk' in 'field list'
                   Insert:   2000     0.00s      0.28 ns/op       0 B/op      0 allocs/op
      MultiInsert 100 row:    500     Not support multi insert
Error 1054: Unknown column 'id,pk' in 'where clause'
                   Update:   2000     0.00s      0.39 ns/op       0 B/op      0 allocs/op
Error 1054: Unknown column 'id,pk' in 'field list'
                     Read:   4000     0.00s      0.11 ns/op       0 B/op      0 allocs/op
Error 1054: Unknown column 'id,pk' in 'field list'
      MultiRead limit 100:   2000     0.00s      0.29 ns/op       0 B/op      0 allocs/op
upper.io
Error 1364: Field 'name' doesn't have a default value
                   Insert:   2000     0.00s      0.37 ns/op       0 B/op      0 allocs/op
      MultiInsert 100 row:    500     Not support multi insert
                   Update:   2000    13.92s      6960274 ns/op    5906 B/op    318 allocs/op
upper: no more rows in this result set
                     Read:   4000     0.00s      0.15 ns/op       0 B/op      0 allocs/op
Error 1364: Field 'name' doesn't have a default value
      MultiRead limit 100:   2000     0.00s      0.30 ns/op       0 B/op      0 allocs/op
dbx
                   Insert:   2000    13.93s      6966821 ns/op    1914 B/op     40 allocs/op
      MultiInsert 100 row:    500    15.27s     30534547 ns/op   69552 B/op    715 allocs/op
                   Update:   2000    12.32s      6158802 ns/op    2606 B/op     59 allocs/op
                     Read:   4000    27.26s      6816054 ns/op    2774 B/op     74 allocs/op
      MultiRead limit 100:   2000    26.33s     13166953 ns/op   78848 B/op   1737 allocs/op
orm
                   Insert:   2000    15.11s      7554026 ns/op    1937 B/op     40 allocs/op
      MultiInsert 100 row:    500    21.92s     43849013 ns/op  147170 B/op   1534 allocs/op
                   Update:   2000    14.95s      7475361 ns/op    1928 B/op     40 allocs/op
                     Read:   4000    30.66s      7664590 ns/op    2800 B/op     97 allocs/op
      MultiRead limit 100:   2000    27.51s     13753703 ns/op   85216 B/op   4287 allocs/op
xorm
                   Insert:   2000    13.89s      6943038 ns/op    2543 B/op     68 allocs/op
      MultiInsert 100 row:    500    17.43s     34853332 ns/op  233982 B/op   4751 allocs/op
                   Update:   2000    13.53s      6765132 ns/op    2800 B/op     96 allocs/op
                     Read:   4000    29.88s      7469699 ns/op    9307 B/op    243 allocs/op
      MultiRead limit 100:   2000    29.00s     14501894 ns/op  180009 B/op   8083 allocs/op
gorm
                   Insert:   2000    26.32s     13160785 ns/op    7336 B/op    149 allocs/op
      MultiInsert 100 row:    500     Not support multi insert
                   Update:   2000    42.39s     21196195 ns/op   19124 B/op    402 allocs/op
                     Read:   4000    28.19s      7047402 ns/op   11611 B/op    239 allocs/op
      MultiRead limit 100:   2000    41.23s     20615999 ns/op  250911 B/op   6225 allocs/op

Reports: 

  2000 times - Insert
       raw:     7.36s      3681234 ns/op     552 B/op     12 allocs/op
      xorm:    13.89s      6943038 ns/op    2543 B/op     68 allocs/op
       dbx:    13.93s      6966821 ns/op    1914 B/op     40 allocs/op
       orm:    15.11s      7554026 ns/op    1937 B/op     40 allocs/op
      hood:    18.30s      9149968 ns/op   12090 B/op    199 allocs/op
      gorm:    26.32s     13160785 ns/op    7336 B/op    149 allocs/op
      gorp:     0.00s      0.28 ns/op       0 B/op      0 allocs/op
  upper.io:     0.00s      0.37 ns/op       0 B/op      0 allocs/op
       qbs:     no primary key field

   500 times - MultiInsert 100 row
       dbx:    15.27s     30534547 ns/op   69552 B/op    715 allocs/op
      xorm:    17.43s     34853332 ns/op  233982 B/op   4751 allocs/op
       raw:    17.59s     35184280 ns/op  110864 B/op    811 allocs/op
       orm:    21.92s     43849013 ns/op  147170 B/op   1534 allocs/op
       qbs:     Not support multi insert
      gorp:     Not support multi insert
  upper.io:     Not support multi insert
      hood:     Not support multi insert
      gorm:     Not support multi insert

  2000 times - Update
       raw:     7.09s      3543659 ns/op     616 B/op     14 allocs/op
       dbx:    12.32s      6158802 ns/op    2606 B/op     59 allocs/op
      xorm:    13.53s      6765132 ns/op    2800 B/op     96 allocs/op
  upper.io:    13.92s      6960274 ns/op    5906 B/op    318 allocs/op
       orm:    14.95s      7475361 ns/op    1928 B/op     40 allocs/op
      hood:    19.83s      9914022 ns/op   12081 B/op    199 allocs/op
      gorm:    42.39s     21196195 ns/op   19124 B/op    402 allocs/op
      gorp:     0.00s      0.39 ns/op       0 B/op      0 allocs/op
       qbs:     no primary key field

  4000 times - Read
       raw:    13.27s      3317756 ns/op    1432 B/op     37 allocs/op
       dbx:    27.26s      6816054 ns/op    2774 B/op     74 allocs/op
      gorm:    28.19s      7047402 ns/op   11611 B/op    239 allocs/op
      xorm:    29.88s      7469699 ns/op    9307 B/op    243 allocs/op
       orm:    30.66s      7664590 ns/op    2800 B/op     97 allocs/op
      hood:    35.64s      8910711 ns/op    4242 B/op     55 allocs/op
      gorp:     0.00s      0.11 ns/op       0 B/op      0 allocs/op
  upper.io:     0.00s      0.15 ns/op       0 B/op      0 allocs/op
       qbs:     no primary key field

  2000 times - MultiRead limit 100
       raw:    20.48s     10241508 ns/op   34704 B/op   1320 allocs/op
       dbx:    26.33s     13166953 ns/op   78848 B/op   1737 allocs/op
       orm:    27.51s     13753703 ns/op   85216 B/op   4287 allocs/op
      xorm:    29.00s     14501894 ns/op  180009 B/op   8083 allocs/op
      hood:    32.98s     16491001 ns/op  232327 B/op   8765 allocs/op
      gorm:    41.23s     20615999 ns/op  250911 B/op   6225 allocs/op
      gorp:     0.00s      0.29 ns/op       0 B/op      0 allocs/op
  upper.io:     0.00s      0.30 ns/op       0 B/op      0 allocs/op
       qbs:     no primary key field

 

```