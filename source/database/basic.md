---
title: Basic Usage
---

## Running SQL Queries

After configuration you can run querys using `goeloquent.Eloquent`

### Select Query

#### Select map
```golang

var row = make(map[string]interface{})
r, err := goeloquent.Eloquent.Select("select * from users where id = ? ", []interface{}{4}, &row)
if err != nil {
    panic(err.Error())
}
fmt.Printf("%#v\n",row)
//map[string]interface {}{"age":0, "balance":100, "email":[]uint8{0x6a, 0x6f, 0x68, 0x6e, 0x40, 0x68, 0x6f, 0x74, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6d}, "id":4, "username":[]uint8{0x6a, 0x6f, 0x68, 0x6e}}
fmt.Printf("%#v\n",r)
//goeloquent.ScanResult{Count:1}

```

#### Select slice of maps
```
var rows []map[string]interface{}
r, err = goeloquent.Eloquent.Select("select * from users limit 10 ", nil, &rows)

```
#### Select struct
```
type BasicUser struct {
	Id       int64
	NickName string
	Age      int
}
var user BasicUser
goeloquent.Eloquent.Select("select * from users where id = ? ", []interface{}{4}, &user)
fmt.Printf("%#v\n", user)
//main.BasicUser{Id:4, NickName:"john", Age:10}
```
#### Select slice of structs
```
var users []BasicUser
goeloquent.Eloquent.Select("select * from users limit 10 ", nil, &users)
fmt.Printf("%#v\n", users)
```

`Select` method first parameter is the sql string, second is the query bindings, third one is the dest.  
For third parameter dest ,you can pass by a pointer of map,a pointer of slice of map , a pointer of struct , a pointer of slice of struct, we will convert it for you .  

Usually every query of this package returns two value. The first is always the Golang standard liberary `database/sql` package `sql.Result`.Second is the raised error.    

When you run insert query , you can call `result.LastInsertId()` to get the created record's id. When you select,update,delete, you can call `result.RowsAffected()` to get the number of affected rows.

### Insert Query
#### Insert Row&Rows
```golang
result, err := goeloquent.Eloquent.Insert("insert into users (nick_name,email,age) values (?,?,?)", []interface{}{"John Doe", "john@hotmail.com", 50})
//result, err := goeloquent.Eloquent.Insert("insert into users (nick_name,email,age) values (?,?,?),(?,?,?),(?,?,?)", []interface{}{"John Doe", "john@hotmail.com", 50, "Jane Doe", "jane@hotmail.com", 50, "Jack", "jack@hotmail.com", 50})

if err != nil {
    panic(err.Error())
}
fmt.Println(result.LastInsertId())
```
### Update Query
#### Update Rows
```golang
result, err := goeloquent.Eloquent.Delete("delete from users where id  = ?", []interface{}{100})
if err != nil {
    panic(err.Error())
}
fmt.Println(result.RowsAffected())
```
### Statement
```golang
_, err := goeloquent.Eloquent.Statement("drop table tests", nil)
	if err != nil {
		panic(err.Error())
	}
```
## Switch Connections
You may have mutiple connections.If you want change the connection or specify use another connection instead of defautl, use `Conn()`, it takes a string as the connection name which you defined in config.
```golang
_, err := goeloquent.Eloquent.Conn("chat").Statement("drop table test", nil)
if err != nil {
    panic(err.Error())
}
```
## Get Raw Connection
In case you need original connection you can use `Raw`,this method return a `*sql.DB`
```
goeloquent.Eloquent.Raw("default")

```
