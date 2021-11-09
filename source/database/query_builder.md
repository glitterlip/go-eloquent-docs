---
title: Query Builder
---

## Introduction

Query builder provides a convenient, fluent interface to creating and running database queries.

## Running Database Queries

While using query builder, you can pass by 4 kinds of pointer as destition.They are

1. pointer of map
2. pointer of slice of map
3. pointer of struct
4. pointer of slice of struct

### Retrieving All Rows From A Table

```
type ChatUser struct {
	Id       int64
	UserName string
	Email    sql.NullString
	Location sql.NullString
}

var users []map[string]interface{}
var userStructs []ChatUser

goeloquent.Eloquent.Conn("chat").Table("users").Get(&users,"id","user_name")
//select `id`,`user_name` from `users` [] {23} 66.375062ms

goeloquent.Eloquent.Conn("chat").Table("users").Get(&userStructs)
//select * from `users` [] {23} 69.783547ms

goeloquent.Eloquent.Conn("chat").Model(&ChatUser{}).Get(&userModelStructs)
//select * from `users` [] {23} 71.613366ms
```

First we define a struct called `ChatUser`,then we use `chat` connection , and get records from table `users`.  
Depends on the type of the first destition parameter,if it's a map , we auto convert it.If it's a struct ,we will convert 
the filed to desired type and assign it.   
Struct field snakename correspond to database table field name.In the
example above,you can find struct `User` field `UserName` turn to `user_name` in sql. If you have a different
mapping,you can add an additional tag to the field. 

```
type ChatUser struct {
	Id       int64
	UserName string `goelo:"column:username"`
	Email    sql.NullString
	Location sql.NullString
}
```

Now we will pair struct field `UserName` with database table field `username`. Like in Laravel ,you can pass by
addtional database columnnames to be selected.

{% note info Note %}
You will find a comment under each query, it's logger print func that we set before.
It's a `goeloquent.Log` struct, include 
1. sql string we just executed
2. sql binding parameters ([]interface{})
3. `sql.Result` with a count of affected/retrived rows
4. `time.Duration` in milliseconds , we use a remote databse to develop,so it include network time.
{% endnote %}

If you use Laravel before , you might know something about `Facade`.
Here we define a variable called `DB` and assign it to `open` function returned resules.

```golang
var DB *goeloquent.DB
func init() {
config := map[string]goeloquent.DBConfig{
    "default": {
        Host:     "127.0.0.1",
        Port:     "3506",
        Database: "eloquent",
        Username: "root",
        Password: "root",
        Driver:   "mysql",
    },
}

DB = goeloquent.Open(config)

```
Now if you use `default` connection , you can just use `DB.Table("users").Get(&userStructs)` to get records! I know , simple ! Right?

### Retrieving A Single Row / Column From A Table

You might just need one single row, use `First`.

```
DB.Table("users").Where("username","john").First(&john)
//{select * from `users` where `username` = ? limit 1 [john] {1} 68.820758ms}
```

If you don't need entire row,just a column value ,use `Value`.
```
var email string
DB.Table("users").Where("username","john").Value(&email,"email")
//{select `email` from `users` where `username` = ? [john] {1} 73.392676ms}
fmt.Println(email)
//john@hotmail.com
```
Use `Find` to get a row by `id`
```golang
var id4 ChatUser
DB.Table("users").Find(&id4,4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 63.692909ms}
fmt.Println(id4)
//{4  {john@hotmail.com true} { false}}

// Even more ,like Laravel
var findMore []ChatUser
DB.Table("users").Find(&findMore,[]interface{}{4,6,8})
//{select * from `users` where `id` in (?,?,?) [4 6 8] {3} 68.072428ms}
fmt.Println(findMore)
//[{4  {john@hotmail.com true} { false}} {6  {john@apple.com true} { false}} {8  {john@yahoo.com true} { false}}]
```
As you can see ,`Find` can take a slice to `Find` mutiple records.

### Retrieving A List Of Column Values
You can use `Pluck` to get a list of a single column
```golang
var titles []string
DB.Table("posts").Pluck(&tites,"title")
//{select `path` from `images` [] {3} 71.442575ms}
fmt.Println(paths)
//[/images1.jpg /images2.jpg /img3.png]
```
## Chunking Results

### Chunk

developing

### ChunkById

developing
## Aggregates
We have `count,max,min,avg,sum` aggregate methods.
```
var total,avg,max,sum,min float64
goeloquent.Eloquent.Conn("default").Model(&DefaultUser{}).Where("age",">",20).Count(&total, "balance")
//{select count(`balance`) as aggregate from `users` where `age` > ? [20] {1} 69.96036ms}

goeloquent.Eloquent.Conn("default").Model(&DefaultUser{}).Max(&max, "balance")
//{select max(`balance`) as aggregate from `users` [] {1} 72.807184ms}

goeloquent.Eloquent.Conn("default").Model(&DefaultUser{}).Min(&min, "balance")
//{select min(`balance`) as aggregate from `users` [] {1} 64.839132ms}

goeloquent.Eloquent.Conn("default").Model(&DefaultUser{}).Sum(&sum, "balance")
//{select sum(`balance`) as aggregate from `users` [] {1} 78.220135ms}

goeloquent.Eloquent.Conn("default").Model(&DefaultUser{}).Avg(&avg, "balance")
//{select avg(`balance`) as aggregate from `users` [] {1} 94.77536ms}

fmt.Println(total)//10
fmt.Println(min)//0
fmt.Println(sum)//1600
fmt.Println(max)//100
fmt.Println(avg)//53.3333
```
## Select Statements
### Specifying A Select Clause
You can use `Select` to specify which column to select

```
var m = make(map[string]interface{})
DB.Conn("chat").Table("users").Select("id","phone","location").First(&m)
//{select `id`,`phone`,`location` from `users` limit 1 [] {1} 62.829325ms}
```
## Raw Expressions
In `Select,Where,GroupBy,OrderBy,Having`,you can pass by a `Expression` . **This could lead to SQL injection **
### Raw Methods
developing
## Joins
developing

## Unions
developing

## Basic Where Clauses

### Where
Usually `Where` function takes 4 parameters,it's column,operator,value,`and/or` Logical Operator.    
Default Logical Operator is `and`,default operator is `=`.

```
var userStructSlice1 []DefaultUser
DB.Table("users").Where("age",">",18,goeloquent.BOOLEAN_AND).Where("balance","=",0,goeloquent.BOOLEAN_AND).Get(&userStructSlice1)
//{select * from `users` where `age` > ? and `balance` = ? [18 0] {1} 64.546997ms}
fmt.Printf("%#v",userStructSlice1)
[]main.DefaultUser{main.DefaultUser{Table:"", Id:6, Age:20, Balance:0, Email:"john@apple.com", NickName:"test3"}}

```

You can skip 4th parameter `and/or` Logical Operator when it's `and`
```
DB.Table("users").Where("age",">",18).Where("balance","=",0).Get(&userStructSlice1)
//select * from `users` where `age` > ? and `balance` = ? [18 0] {2} 62.012133ms
```

You can skip 2nd parameter when it's `=`
```
DB.Table("users").Where("age",">",18).Where("balance",0).Get(&userStructSlice1)
//select * from `users` where `age` > ? and `balance` = ? [18 0] {2} 62.202419ms
```

You can pass by a `[][]interface{}`,each element should be a `[]interface` containing the four parameters that pass to the `where` function  

```
DB.Table("users").Where([][]interface{}{
    {"age", ">", 18, goeloquent.BOOLEAN_AND},
    {"balance", "=", 0, goeloquent.BOOLEAN_AND},
}).Get(&userStructSlice1)
//select * from `users` where `age` > ? and `balance` = ? [18 0] {2} 62.523877ms
```
skip parameters works too
```
DB.Table("users").Where([][]interface{}{
    {"age", ">", 18},
    {"balance", 0},
}).Get(&userStructSlice1)
//select * from `users` where `age` > ? and `balance` = ? [18 0] {2} 61.789099ms
```
## Or Where Clauses
For more readable reason,you may want a `OrWhere` function
```
DB.Table("users").Where("age",">",18).OrWhere("balance","=",0).Get(&userStructSlice1)
select * from `users` where `age` > ? or `balance` = ? [18 0] {24} 62.61687ms
```

## Additional Where Clauses
### WhereBetween/OrWhereBetween/WhereNotBetween/OrWhereNotBetween
```
DB.Table("users").Where("balance", ">", 100).WhereBetween("age", []interface{}{18, 35}).Get(&userStructSlice)
//select * from `users` where `balance` > ? and `age` between ? and ? [100 18 35] {0} 68.290583ms

DB.Table("users").WhereNotBetween("age", []interface{}{18, 35}).Get(&userStructSlice)
//select * from `users` where `age` not between ? and ? [18 35] {23} 69.032302ms

DB.Table("users").Where("balance", ">", 100).OrWhereNotBetween("age", []interface{}{18, 35}).Get(&userStructSlice)
//select * from `users` where `balance` > ? or `age` not between ? and ? [100 18 35] {23} 62.927148ms

DB.Table("users").Where("balance", ">", 100).OrWhereBetween("age", []interface{}{18, 35}).Get(&userStructSlice)
//select * from `users` where `balance` > ? or `age` between ? and ? [100 18 35] {7} 63.241122ms
```
### WhereIn/OrWhereIn/WhereNotIn/OrWhereNotIn
```
goeloquent.Eloquent.Model(&DefaultUser{}).WhereIn("id", []interface{}{1,2,3}).Get(&userStructSlice)
//select * from `users` where `id` in (?,?,?) [1 2 3] {1} 62.159353ms

goeloquent.Eloquent.Model(&DefaultUser{}).WhereNotIn("id", []interface{}{2,3,4}).Get(&userStructSlice)
//select * from `users` where `id` not in (?,?,?) [2 3 4] {28} 68.078067ms

DB.Table("users").Where("username","john").OrWhereIn("email", []interface{}{"john@gmail.com","john@hotmail.com","john@apple.com","john@outlook.com"}).Get(&userStructSlice)
//select * from `users` where `username` = ? or `email` in (?,?,?,?) [john john@gmail.com john@hotmail.com john@apple.com john@outlook.com] {3} 61.692218ms

DB.Table("users").Where("username","joe").OrWhereNotIn("email", []interface{}{"joe@gmail.com","joe@hotmail.com","joe@apple.com","joe@outlook.com"}).Get(&userStructSlice)
//select * from `users` where `username` = ? or `email` not in (?,?,?,?) [joe joe@gmail.com joe@hotmail.com joe@apple.com joe@outlook.com] {30} 64.416506ms
```
### WhereNull/OrWhereNull/OrWhereNotNull/WhereNotNull
```
DB.Table("users").WhereIn("id", []interface{}{1,2,3}).WhereNull("email").Get(&userStructSlice)
//select * from `users` where `id` in (?,?,?) and `email` is null  [1 2 3] {0} 61.984595ms

DB.Table("users").WhereNotIn("id", []interface{}{2,3,4}).WhereNotNull("email").Get(&userStructSlice)
//select * from `users` where `id` not in (?,?,?) and `email` is  not null  [2 3 4] {27} 62.228735ms

DB.Table("users").Where("username","john").OrWhereNull("email").Get(&userStructSlice)
//select * from `users` where `username` = ? or `email` is null  [john] {1} 62.454664ms

DB.Table("users").Where("username","joe").OrWhereNotNull("email").Get(&userStructSlice)
//select * from `users` where `username` = ? or `email` is  not null  [joe] {29} 62.256084ms
```
### WhereDate/WhereMonth/WhereDay/WhereYear/WhereTime
```
var now = time.Now()
fmt.Println(now)
//2021-11-03 16:00:35.461691 +0800 CST m=+0.166644409
DB.Table("users").WhereDate("created_at", now).Get(&userStructSlice)
//{select * from `users` where date(`created_at`) = ? [2021-11-03] {0} 65.800819ms}

DB.Table("users").WhereDate("created_at", "2008-01-03").Get(&userStructSlice)
//{select * from `users` where date(`created_at`) = ? [2008-01-03] {0} 66.675012ms}

DB.Table("users").WhereDay("created_at", now).Get(&userStructSlice)
//{select * from `users` where day(`created_at`) = ? [03] {0} 65.159437ms}

DB.Table("users").WhereDay("created_at", "06").Get(&userStructSlice)
//{select * from `users` where day(`created_at`) = ? [06] {0} 64.92847ms}

DB.Table("users").WhereMonth("created_at", now).Get(&userStructSlice)
//{select * from `users` where month(`created_at`) = ? [11] {10} 70.454652ms}

DB.Table("users").WhereMonth("created_at", "06").Get(&userStructSlice)
//{select * from `users` where month(`created_at`) = ? [11] {10} 66.694005ms}

DB.Table("users").WhereYear("created_at", now).Get(&userStructSlice)
//{select * from `users` where year(`created_at`) = ? [2021] {11} 64.805563ms}

DB.Table("users").WhereYear("created_at", "2020").Get(&userStructSlice)
//{select * from `users` where year(`created_at`) = ? [2020] {0} 64.970053ms}

DB.Table("users").WhereTime("created_at", now).Get(&userStructSlice)
//{select * from `users` where time(`created_at`) = ? [16:00:35] {0} 65.73327ms}

DB.Table("users").WhereTime("created_at", "3:05:16").Get(&userStructSlice)
//{select * from `users` where time(`created_at`) = ? [3:05:16] {0} 66.24917ms}
```
### WhereColumn/OrWhereColumn
```
DB.Table("users").WhereColumn("age", "=", "balance").Get(&userStructSlice)
//{select * from `users` where `age` = `balance` [] {1} 65.095414ms}

DB.Table("users").Where("id",4).OrWhereColumn("age", "=", "balance").Get(&userStructSlice)
//{select * from `users` where `id` = ? or `age` = `balance` [4] {2} 66.101059ms}
```
## Logical Grouping
If you need to group an `where` condition within parentheses,you can pass by a function to  `Where` or use `WhereNested` function
```
DB.Table("users").Where("age", ">", 30).OrWhere(func(builder *goeloquent.Builder) {
    builder.Where("age", ">", 18)
    builder.Where("balance", ">", 5000)
}).Get(&userStructSlice, "username", "email")
//select `username`,`email` from `users` where `age` > ? or (`age` > ? and `balance` > ?) [30 18 5000] {8} 62.204423ms

DB.Table("users").Where("age", ">", 30).WhereNested([][]interface{}{
    {"age", ">", 18},
    {"balance", ">", 5000},
},goeloquent.BOOLEAN_OR).Get(&userStructSlice, "username", "email")
//select `username`,`email` from `users` where `age` > ? or (`age` > ? and `balance` > ?) [30 18 5000] {8} 64.868523ms	
```
## Subquery Where Clauses
```
DB.Table("users").Where("age", ">", 0).WhereSub("id","in", func(builder *goeloquent.Builder) {
    builder.From("users").Where("balance",">",0).Select("id")
    //don't use any finisher function like first/find/get/value/pluck,otherwise it will turn into execute two seperated sql 
},goeloquent.BOOLEAN_OR).Get(&userStructSlice)
```
## Conditional Clauses
```
DB.Model(&models.User{}).When(false, func(builder *goeloquent.Builder) {
    q.Where("id",10)
}).Get(&us)
```
## Ordering, Grouping, Limit & Offset
### Ordering
default order is `asc`
```
DB.Table("users").Where("age", ">", 0).OrderBy("balance",goeloquent.ORDER_DESC).OrderBy("id").Get(&userStructSlice)
//{select * from `users` where `age` > ? order by `balance` desc , `id` asc [0] {24} 73.264891ms}
```
### Grouping
```
var m []map[string]interface{}
goeloquent.Eloquent.Table("comments").GroupBy("commentable_type").Get(&m,"commentable_type",goeloquent.Expression{Value: "count(*) as c"})
//{select `commentable_type`,count(*) as c from `comments` group by `commentable_type` [] {2} 64.624213ms}
fmt.Println(string(m[0]["commentable_type"].([]byte)))
//post
fmt.Println(m[0]["c"])
//2
```
When using `Select` to specify columns, we will quote columns with "`" as string, if you want avoid this, use `Expression`

### Having
```
var m []map[string]interface{}
goeloquent.Eloquent.Table("comments").GroupBy("commentable_type").Having("c",">",2).Get(&m,"commentable_type",goeloquent.Expression{Value: "count(*) as c"})
//{select `commentable_type`,count(*) as c from `comments` group by `commentable_type` having `c` > ? [2] {1} 66.393615ms}
fmt.Println(string(m[0]["commentable_type"].([]byte)))
//video
fmt.Println(m[0]["c"])
//3
```
#### HavingBetween
```
var m []map[string]interface{}
goeloquent.Eloquent.Table("comments").GroupBy("commentable_type").HavingBetween("c",[]interface{}{0,3}).Get(&m,"commentable_type",goeloquent.Expression{Value: "count(*) as c"})
//{select `commentable_type`,count(*) as c from `comments` group by `commentable_type` having `c` between? and ?  [0 3] {2} 65.1953ms}
fmt.Println(string(m[0]["commentable_type"].([]byte)))
//video
fmt.Println(m[0]["c"])
//3
```
### Limit & Offset
```
DB.Conn("chat").Model(&ChatPkUser{}).Limit(5).Offset(3).Select("id","phone","location").Get(&userStructSlice)
//{select `id`,`phone`,`location` from `users` limit 5 offset 3 [] {5} 76.978184ms}
```
## Insert Statement
### Insert map
```
var userMap = map[string]interface{}{
    "username": fmt.Sprintf("%s%d", "Only", time.Now().Unix()),
    "balance":  100,
    "age":      50,
}
result, err := DB.Table("users").Only("balance", "username").Insert(&userMap)
//{insert into `users` ( `username`,`balance` ) values  ( ? , ? )  [Only1635947804 100] {0xc00007a000 0xc00001e160} 78.784832ms}
if err != nil {
    panic(err.Error())
}
insertId, _ := result.LastInsertId()
var inserted = make(map[string]interface{})
DB.Table("users").Where("id",insertId).First(&inserted)
//{select * from `users` where `id` = ? limit 1 [162] {1} 67.92676ms}
fmt.Println(inserted)
//map[age:0 balance:100 id:162 username:[79 110 108 121 49 54 51 53 57 52 55 56 48 52]]
```

While updating/inserting,you can use `Only` to specify which columns to include ,you can use `Except` to specify which columns to exclude
### Batch Insert Map
You can pass by a slice of map to insert several records at once
```
s := []map[string]interface{}{
    {
        "id":       1,
        "username": "userr1",
        "balance":  1000,
        "age":      20,
    },
    {
        "username": "userr2",
        "balance":  50000,
        "age":      50,
    },
}
result, err := DB.Table("users").Except("id").Insert(&s)
//{insert into `users` ( `age`,`username`,`balance` ) values  ( ? , ? , ? )  ,  ( ? , ? , ? )  [20 userr1 1000 50 userr2 50000] {0xc00010c000 0xc00001e160} 86.2694ms}
if err != nil {
    panic(err.Error())
}
fmt.Println(result.LastInsertId())
//168 <nil> 
fmt.Println(result.RowsAffected())
//2 <nil> 
```
When batch insert,`LastInsertId()` will return the first record id
### Insert Struct
```
type Post struct {
    Table   string `goelo:"TableName:posts"`
    Id      int64  `goelo:"primaryKey"`
    UserId  int64
    Title   string
    Summary string
    Content string
}
fmt.Println("table struct insert Only")
var post = Post{
    Id:      10,
    UserId:  2,
    Title:   fmt.Sprintf("%s%d", "table struct insert Only", time.Now().Unix()+1),
    Summary: "Summary",
    Content: fmt.Sprintf("%s%d", "Summary table struct insert Only", time.Now().Unix()+2),
}
//{insert into `posts` ( `summary`,`content`,`user_id`,`title` ) values  ( ? , ? , ? , ? )  [Summary Summary table struct insert Only1635949167 2 table struct insert Only1635949166] {0xc0000da120 0xc0000b0540} 68.205202ms}
result, err := DB.Table("posts").Only("user_id", "title", "content", "summary").Insert(&post)
if err != nil {
    panic(err.Error())
}
fmt.Printf("%#v",post)
//main.Post{Table:"", Id:174, UserId:2, Title:"table struct insert Only1635949166", Summary:"Summary", Content:"Summary table struct insert Only1635949167"}
fmt.Println(result.LastInsertId())
//174 <nil>
```
If you add an tag ``goelo:"primaryKey"`` on primaryKey field , we will update it for you ,otherwise it is its original value

### Batch Insert Structs
```
s := []Post{
    {
        Id:      10,
        UserId:  4,
        Title:   fmt.Sprintf("%s%d", "table slice struct insert ", time.Now().Unix()+1),
        Summary: "Summary",
        Content: fmt.Sprintf("%s%d", "table slice struct insert ", time.Now().Unix()+2),
    },
    {
        Id:      10,
        UserId:  4,
        Title:   fmt.Sprintf("%s%d", "table slice struct insert ", time.Now().Unix()+1),
        Summary: "Summary",
        Content: fmt.Sprintf("%s%d", "table slice struct insert ", time.Now().Unix()+2),
    },
}
result, err = DB.Table("posts").Except("id").Insert(&s)
if err != nil {
    panic(err.Error())
}
```
Insert can accept a pointer of slice of struct
## Update Statements
```
r, err := DB.Table("users").Where([][]interface{}{
    {"age", 18},
    {"balance", 0},
}).Update(map[string]interface{}{
    "balance": 100,
})
//{update `users` set `balance` = ? where `age` = ? and `balance` = ? [100 18 0] {0xc000204000 0xc0002260d0} 80.266387ms}
if err != nil {
    panic(err.Error())
}
fmt.Println(r.RowsAffected())
//2 <nil>
```
## Use Expression
```
r, err := DB.Table("users").Where([][]interface{}{
    {"age", 18},
    {"balance", "!=", 0},
}).Update(map[string]interface{}{
    "balance": goeloquent.Expression{Value: "balance + 100"},
})
//{update `users` set `balance` = balance + 100 where `age` = ? and `balance` != ? [18 0] {0xc0000de120 0xc0000b83a0} 75.251657ms}
if err != nil {
    panic(err.Error())
}
fmt.Println(r.RowsAffected())
//2 <nil>
```
Another Example
```
r, err := DB.Table("users").Update(map[string]interface{}{
    "balance": goeloquent.Expression{Value: "balance + age*100"},
})
//{update `users` set `balance` = balance + age*100 [] {0xc000110000 0xc000128090} 68.057742ms}
if err != nil {
    panic(err.Error())
}
fmt.Println(r.RowsAffected())
//32 <nil>
```
You can use `Expression` to update column base on another column
## Delete
```
r,err:=DB.Table("users").Where("id",2).Delete()
if err!=nil {
    panic(err.Error())
}
fmt.Println(r.RowsAffected())
//{ delete from `users` where `id` = ? [2] {0xc000102000 0xc00001e150} 73.972219ms}
//1 <nil>
```
## Pessimistic Locking

```
var us []map[string]interface{}
DB.Table("users").LockForUpdate().Where("id", "<", 100).Get(&us)
//{select * from `users` where `id` < ? for update  [100] {16} 66.21529ms}

DB.Table("users").SharedLock().Where("id", "<", 100).Get(&us)
//{select * from `users` where `id` < ? lock in share mode  [100] {16} 67.434753ms}
```
