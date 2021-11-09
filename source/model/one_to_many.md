---
title: One To Many Relation
---

## One To Many

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
func (u *User) TableName() string {
	return "users"
}
func (u *User) PostsRelation() *goeloquent.RelationBuilder {
   // relatedParentKey (post.user_id) equals to parentKey (user.id)
	return u.HasMany(u, &Post{}, "user_id", "id")
}

```

1. add an tag ``goelo:"HasMany:RelationMethod"`` for one to one many field type can be slice of struct,slice of pointer
   of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `HasMany`
3. `HasMany` takes 4 parameter,first one is a pointer of current model,second is a pointer of related model, third is
   related model database field correspond to parent(current) model,last one is parent(current) model field correspond
   to third parameter.

### Usage Example

Use With when retrive

```
var u models.User
_, _ = goeloquent.Eloquent.Model(&models.User{}).With("Posts").Find(&u, 4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 60.14015ms}
//{select * from `posts` where `user_id` is  not null  and `user_id` in (?) [4] {2} 59.815928ms}

var users []models.User
_,_ = DB.Model(&models.User{}).With("Posts").Get(&users)
//{select * from `users` where `id` < ? [10] {3} 60.340493ms}
//{select * from `posts` where `user_id` is  not null  and `user_id` in (?,?,?) [4 6 8] {2} 60.569247ms}

```

Use Load Or Directly Call Relation Method

```
var uu models.User
DB.Model(&models.User{}).Find(&uu,4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 60.349539ms}

uu.Load("Posts")
//{select * from `posts` where `user_id` is  not null  and `user_id` in (?) [4] {2} 60.349589ms}

var posts []models.Post
uu.PostsRelation().Get(&posts)
//{select * from `posts` where `user_id` = ? and `user_id` is  not null  [4] {2} 59.872067ms}
```

## One To Many (Inverse) / Belongs To

```
type Post struct {
	goeloquent.EloquentModel
	Id         int64      `goelo:"column:id;primaryKey"`
	UserId     int64      `goelo:"column:user_id"`
	Title      string     `goelo:"column:title"`
	Content    string     `goelo:"column:content"`
    User       *User      `goelo:"BelongsTo:UserRelation"`
}
func (p *Post) TableName() string {
	return "posts"
}
func (p *Post) UserRelation() *goeloquent.RelationBuilder {
	return p.BelongsTo(p,&User{},"user_id","id")
}
```

1. add an tag ``goelo:"BelongsTo:RelationMethod"`` for one to one many field type can be struct, pointer of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `BelongsTo`
3. `BelongsTo` takes 4 parameter,first one is a pointer of current model,second is a pointer of related model, third is
   parent(current) model database field correspond to related model,last one is related model field correspond to third
   parameter.

### Usage Example

Use With when retrive

```
var ps []models.Post
DB.Model(&models.Post{}).With("User").Get(&ps)
//{select * from `posts` [] {6} 61.408829ms}
//{select * from `users` where `id` in (?,?,?,?,?,?) [4 2 2 2 2 4] {1} 60.460676ms}

```

Use Load Or Directly Call Relation Method

```
ps[0].Load("User")
var user models.User
//{select * from `users` where `id` in (?) [4] {1} 62.15223ms}
ps[0].UserRelation().Get(&user)
//{select * from `users` where `id` = ? [4] {1} 60.82013ms}
```
