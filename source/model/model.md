---
title: Model
---

Model is a more advanced struct.With **Relationships**,**Events**,**Connection**,**TimeStamps** features.

## Conventions

First,let's see how to define a model.

1. A model should be a struct with an embed anonymous struct of type `goeloquent.EloquentModel`
2. primaryKey should be `int64` with tag `goelo:"primaryKey"`
3. **Each database field** should have a tag `goelo:"column:databasefiledname"`
4. By default , we will use the model's SnakeCase name as table name (e.g. 'User' as 'user','UserInfo' as 'user_info'')
   ,you can change this by implement `type TableName interface`
5. By default , we will use the connection with name `default`.You can change this by
   implement `type ConnectionName interface`.
6. When creating model,we will set timestamps for you if you have a tag `goelo:"column:created_at;CREATED_AT"`
6. When updating model,we will set timestamps for you if you have a tag `goelo:"column:updated_at;UPDATED_AT"`

Example of User model

```
type UserModel struct {
	goeloquent.EloquentModel
	Id        int64          `goelo:"column:id;primaryKey"`
	UserName  string         `goelo:"column:user_name"`
	NickName  sql.NullString `goelo:"column:nick_name"`
	Phone     sql.NullString `goelo:"column:phone"`
	Age       int            `goelo:"column:age"`
	CreatedAt sql.NullTime   `goelo:"column:created_at;CREATED_AT"`
	UpdatedAt sql.NullTime   `goelo:"column:updated_at;UPDATED_AT"`
}

func (u *UserModel) TableName() string {
	return "users"
}
func (u *UserModel) ConnectionName() string {
	return "chat"
}
```
## Model Events
```
func (u *UserModel) Creating(user *UserModel) bool {
	if user.Age <= 0 {
		user.Age = 10 //change it use `user` instead of `u`
	}
	fmt.Println(user.Exists)//false
	user.GetOrigin()//map with default values
	user.GetDirty() //map with non-zero values
	user.GetChanges() //it's not saved to dasebase so empty
	if strings.Contains(user.NickName.String, "shit") {
	     // bad name - -
		return false //return false to abort create
	}
	return true
}
func (u *UserModel) Created(user *UserModel) {
	user.GetOrigin() //same with creating
	user.GetDirty() //same with creating and an `id` value
	user.GetChanges() //same with creating
}
func (u *UserModel) Updating(user *UserModel) bool {
	user.GetOrigin()//original databse value
	user.GetDirty()//map with changed values
	user.GetChanges() // empty map
	return true
}
func (u *UserModel) Updated(user *UserModel) {
	user.GetOrigin()//same with updating
	user.GetDirty()//same with updating
	user.GetChanges()//same with GetDirty
}
func (u *UserModel) Saving(user *UserModel) bool {
	user.GetOrigin()//when create same with creating,when update same with updating
	user.GetDirty()//when create same with creating,when update same with updating
	user.GetChanges()//when create same with creating,when update same with updating
	return true
}
func (u *UserModel) Saved(user *UserModel) {
	user.GetOrigin() //when create same with created
	user.GetDirty()//when create same with created
	user.GetChanges()//when create same with created
}
func (u *UserModel) Deleting(user *UserModel) {
	user.GetOrigin()
	user.GetDirty()
	user.GetChanges()
}
func (u *UserModel) Deleted(user *UserModel) {
	user.GetOrigin()
	user.GetDirty()
	user.GetChanges()
}

```

## Retrive Models

```
var us []UserModel
DB.Model(&UserModel{}).Where("age", ">", 0).Get(&us)
//select * from `user` where `age` > ? [0] {23} 61.963146ms
fmt.Println(us)

var u UserModel
DB.Model(&UserModel{}).Find(&u,22)
//{select * from `users` where `age` > ? [0] {23} 76.868039ms}
```

## Create Model

```
var nu UserModel
goeloquent.NewModel(&nu)
nu.Age = 23
nu.UserName = "hello"
nu.Phone = sql.NullString{
  String: "+19785920200",
  Valid:  true,
}
nu.Create()
//{insert into `users` ( `user_name`,`phone`,`age`,`created_at` ) values  ( ? , ? , ? , ? )  [hello {+19785920200 true} 23 {2021-11-05 16:21:47.188486 +0800 CST m=+0.570200737 true}] <nil> 65.652297ms}
```

## Update Model

```
var u UserModel
DB.Model(&UserModel{}).Find(&u, 22)
u.Age = 23
u.Save()
//{update `users` set `user_name` = ? , `nick_name` = ? , `phone` = ? , `age` = ? , `created_at` = ? , `updated_at` = ? 
//where `id` = ? [s {s true} {1965320233 true} 23 {2021-04-07 15:33:21 +0000 UTC true} {2021-11-05 16:13:19.467342 +0800 CST m=+0.424807671 true} 22] {0xc000162360 0xc00012d110} 64.493242ms}
```

## Mass Update

```
var u UserModel
DB.Model(&u).Where("age", ">", 200).Update(map[string]interface{}{
  "age":   goeloquent.Expression{Value: "age + 10"},
  "phone": "unknown",
})
//{update `users` set `age` = age + 10 , `phone` = ? where `age` > ? [unknown 200] {0xc000140090 0xc00001eea0} 72.325303ms}
```

You can use `Only/Except` to specify which columns to be inserted/updated

## Scopes

developing

## Events

We have `saving,saved,creating,created,updating,updated,deleting,deleted` model events.

### Create

When Creating,we will call event methods in this order.

```
var nu UserModel
goeloquent.NewModel(&nu)
nu.Age = 23
nu.UserName = "hello"
nu.Phone = sql.NullString{
  String: "+19785920200",
  Valid:  true,
}
nu.Create()
//call `Saving` method. if it return false, abort sql execution with err "abort by EventSaving func"
//call `Creating` method. if it return false, abort sql execution with err "abort by EventSaving func"
//execute sql.if any error raises return it.
//call `Created` method
//call `Saved` method
```

### Update

```
var u UserModel
DB.Model(&UserModel{}).Find(&u, 22)
u.Age = 23
u.Save()
//call `Saving` method. if it return false, abort sql execution with err "abort by EventSaving func"
//call `Updating` method. if it return false, abort sql execution with err "abort by EventSaving func"
//execute sql.if any error raises return it.
//call `Updated` method
//call `Saved` method
```

### Delete

```
var u UserModel
DB.Model(&UserModel{}).Find(&u, 22)
u.Delete()
//call `Deleting` method. if it return false, abort sql execution with err "abort by EventSaving func"
//execute sql.if any error raises return it.
//call `Deleted` method
```

## GetOrigin/GetDirty/GetChange
You can call `GetOrigin` to get the original attributes since you retrive it from database.  
If it's just created and haven't been save to database it will be a `map[string]interface` with all default values.
You can get a `map[string]interface` with all changed value than hasn't been save to databse yet by call `GetDirty`.
{% note warn warn %}
Inside event method,we use `reflect` package, so if you want change model's field value ,don't use method pointer receiver,instead use method parameter(which is also a same type model pointer)
{% endnote %}
### Mute Events
If you don't want trigger event use `Mute()`

```
var u UserModel
DB.Model(&UserModel{}).Find(&u, 22)
u.Age = 23
u.Mute(goeloquent.EventUpdating,goeloquent.EventUpdated,goeloquent.EventSaving,goeloquent.EventSaved).Save()

```
