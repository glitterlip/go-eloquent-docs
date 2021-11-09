---
title: One To One Relation
---
## One To One

```
type User struct {
	goeloquent.EloquentModel
	Id       int64          `goelo:"column:id;primaryKey"`
	UserName sql.NullString `goelo:"column:username"`
	Age      int            `goelo:"column:age"`
	Balance  int            `goelo:"column:balance"`
	Email    sql.NullString `goelo:"column:email"`
	Info     Info           `goelo:"HasOne:InfoRelation"`
    Posts    []Post         `goelo:"HasMany:PostsRelation"`

}
func (u *User) InfoRelation() *goeloquent.RelationBuilder {
    //this means parentKey (user.id) equals to relatedParentKey (info.user_id)
	return u.HasOne(u, &Info{}, "user_id", "id")
}

func (u *User) TableName() string {
	return "users"
}
```

1. add an tag ``goelo:"HasOne:RelationMethod"`` for one to one model field type can be struct,pointer of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `HasOne`
3. `HasOne` takes 4 parameter,first one is a pointer of current model,second is a pointer of related model, third is
related model database field correspond to current(parent) model,last one is current model field correspond to third parameter.

### Usage Example

Use With when retrive

```
var u models.User
_, _ = DB.Model(&models.User{}).With("Info").Find(&u, 4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 68.017982ms}
//{select * from `infos` where `user_id` is  not null  and `user_id` in (?) [4] {1} 66.096682ms}


var user []models.User
_, _ = DB.Model(&models.User{}).Where("id", "<", 10).With("Info").Get(&user)
//{select * from `users` where `id` < ? [10] {3} 60.846906ms}
//{select * from `infos` where `user_id` is  not null  and `user_id` in (?,?,?) [4 6 8] {1} 60.688878ms}
```

Use Load Or Directly Call Relation Method

```
var e models.User
DB.Model(&e).Find(&e,4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 59.984896ms}
e.Load("Info")
//{select * from `infos` where `user_id` is  not null  and `user_id` in (?) [4] {1} 59.832026ms}

var info models.Info
e.InfoRelation().Get(&info)
//{select * from `infos` where `user_id` = ? and `user_id` is  not null  [4] {1} 59.928416ms}
```

### Defining The Inverse Of The Relationship

```
type Info struct {
	goeloquent.EloquentModel
	Id       int64     `goelo:"column:id;primaryKey"`
	UserId   int64     `goelo:"column:user_id"`
	Address  string    `goelo:"column:address"`
	Gender   int       `goelo:"column:gender"`
	Phone    string    `goelo:"column:phone"`
	Birthday time.Time `goelo:"column:birthday"`
	User     *User     `goelo:"BelongsTo:UserRelation"`

}
func (i *Info) TableName() string {
	return "infos"
}
func (i *Info) UserRelation() *goeloquent.RelationBuilder {
    //this means parentRelatedKey (info.user_id) equals to relatedKey (user.id)
	return i.BelongsTo(i, &User{}, "user_id", "id")
}
```

1. add an tag ``goelo:"BelongsTo:RelationMethod"`` for one to one model field type can be struct,pointer of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `BelongsTo`
3. `BelongsTo` takes 4 parameter,first one is a pointer of current model,second is a pointer of related model, third is
parent(current) model database field correspond to related model,last one is related model field correspond to third
parameter.

### Usage Example

Use With when retrive

```
var u Info
_, err := DB.Model(&UserInfo{}).Where("id", 2).With("User").First(&u)
//{select * from `infos` where `id` = ? limit 1 [2] {1} 63.14363ms}
//{select * from `users` where `id` in (?) [4] {1} 59.791525ms}

var us []Info
_, err = DB.Model(&UserInfo{}).With("User").Get(&us)
//{select * from `infos` [] {2} 59.996246ms}
//{select * from `users` where `id` in (?,?) [4 2] {1} 62.111639ms}
```

Use Load Or Directly Call Relation Method

```
var info models.Info
_, _ = DB.Model(&models.Info{}).With("User").Find(&info, 4)
//{select * from `infos` where `id` in (?) limit 1 [4] {1} 60.777364ms}
//{select * from `users` where `id` in (?) [2] {0} 61.371022ms}

var infos []models.Info
_, _ = DB.Model(&models.Info{}).Where("id", "<", 10).With("User").Get(&infos)
//{select * from `infos` where `id` < ? [10] {4} 60.135982ms}
//{select * from `users` where `id` in (?,?,?,?) [4 2 8 6] {3} 61.204417ms}


var i models.Info
DB.Model(&i).Find(&i, 4)
//{select * from `infos` where `id` in (?) limit 1 [4] {1} 59.965612ms}
i.Load("User")
//{select * from `users` where `id` in (?) [2] {0} 59.522169ms}

var user models.User
i.UserRelation().Get(&user)
//{select * from `users` where `id` = ? [2] {0} 60.132747ms}
```

{% note warn HasOne/BelongsTo %} If User has a field Info and Info has a field User,that's not gonna work. IDE will warn
you `Invalid recursive type 'Info' Info → User → Info`,so define one of them as pointer

{% endnote %}
