---
title: Many To Many
---

## Many To Many

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
	Roles    []Role         `goelo:"BelongsToMany:RolesRelation"`

}
func (u *User) TableName() string {
	return "users"
}
func (u *User) RolesRelation() *goeloquent.RelationBuilder {
      //pivotTable.pivotParentKey(user_roles.user_id) = parentTable.parentKey (user.id),pivotTable.pivotRelatedKey(user_roles.role_id) = relatedTabel.relatedKey (roles.id)
	return u.BelongsToMany(u, &Role{}, "user_roles", "user_id","role_id","id","id")
}

```

1. add an tag ``goelo:"BelongsToMany:RelationMethod"`` for many to one many field type can be slice of struct,slice of
   pointer of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `BelongsToMany`
3. `BelongsToMany` takes 6 parameter,first one is a pointer of current model,second is a pointer of related model, third
   is pivotTable , 4th is pivotTable.pivotParentKey , 5th is pivotTable.pivotRelatedKey ,last one is related model field
   correspond to 4th parameter.

### Usage Example

Use With when retrive

```
var u models.User
var us []models.User
DB.Model(&models.User{}).With("Roles").Find(&u,4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 61.560446ms}
//{select `roles`.*,`user_roles`.`user_id` as `goelo_pivot_user_id`,`user_roles`.`role_id` as `goelo_pivot_role_id` from `roles` inner join user_roles on `user_roles`.`role_id` = `roles`.`id` where `user_roles`.`user_id` in (?) [4] {2} 62.289392ms}

DB.Model(&models.User{}).With("Roles").Where("id", "<", 10).Get(&us)
//{select * from `users` where `id` < ? [10] {3} 60.056161ms}
//{select `roles`.*,`user_roles`.`user_id` as `goelo_pivot_user_id`,`user_roles`.`role_id` as `goelo_pivot_role_id` from `roles` inner join user_roles on `user_roles`.`role_id` = `roles`.`id` where `user_roles`.`user_id` in (?,?,?) [4 6 8] {4} 60.250975ms}


```

Use Load Or Directly Call Relation Method

```
var uu models.User
var roles []models.Role
DB.Model(&uu).Find(&uu,6)
//{select * from `users` where `id` in (?) limit 1 [6] {1} 60.646603ms}

uu.Load("Roles")
//{select `roles`.*,`user_roles`.`user_id` as `goelo_pivot_user_id`,`user_roles`.`role_id` as `goelo_pivot_role_id` from `roles` inner join user_roles on `user_roles`.`role_id` = `roles`.`id` where `user_roles`.`user_id` in (?) [6] {2} 59.965972ms}

uu.RolesRelation().Get(&roles)
//{select `roles`.*,`user_roles`.`user_id` as `goelo_pivot_user_id`,`user_roles`.`role_id` as `goelo_pivot_role_id` from `roles` inner join user_roles on `user_roles`.`role_id` = `roles`.`id` where `user_id` = ? [6] {2} 59.866877ms}
```

## Many To Many (Inverse)

```
type Role struct {
	goeloquent.EloquentModel
	Id          int64  `goelo:"column:id;primaryKey"`
	Name        string `goelo:"column:name"`
	DisplayName string `goelo:"column:display_name"`
	Users       []User `goelo:"BelongsToMany:UsersRelation"`
}

func (r *Role) UsersRelation() *goeloquent.RelationBuilder {
	return r.BelongsToMany(r, &User{}, "user_roles", "role_id", "user_id", "id", "id")
}
func (r *Role) TableName() string {
	return "roles"
}
```

Well same to BelongsToMany

### Usage Example

Use With when retrive

```
var r models.Role
var rs []models.Role
DB.Model(&models.Role{}).With("Users").Find(&r,10)
//{select * from `roles` where `id` in (?) limit 1 [10] {1} 59.997077ms}
//{select `users`.*,`user_roles`.`role_id` as `goelo_pivot_role_id`,`user_roles`.`user_id` as `goelo_pivot_user_id` from `users` inner join user_roles on `user_roles`.`user_id` = `users`.`id` where `user_roles`.`role_id` in (?) [10] {2} 61.439653ms}

DB.Model(&models.Role{}).With("Users").Get(&rs)
//{select * from `roles` [] {4} 60.020418ms}
//{select `users`.*,`user_roles`.`role_id` as `goelo_pivot_role_id`,`user_roles`.`user_id` as `goelo_pivot_user_id` from `users` inner join user_roles on `user_roles`.`user_id` = `users`.`id` where `user_roles`.`role_id` in (?,?,?,?) [10 12 14 16] {4} 61.215167ms}

```

Use Load Or Directly Call Relation Method

```
var rr models.Role
var users []models.User
DB.Model(&rr).Find(&rr,10)
//{select * from `roles` where `id` in (?) limit 1 [10] {1} 62.39588ms}
rr.Load("Users")
//{select `users`.*,`user_roles`.`role_id` as `goelo_pivot_role_id`,`user_roles`.`user_id` as `goelo_pivot_user_id` from `users` inner join user_roles on `user_roles`.`user_id` = `users`.`id` where `user_roles`.`role_id` in (?) [10] {2} 60.394408ms}

rr.UsersRelation().Get(&users)
//{select `users`.*,`user_roles`.`role_id` as `goelo_pivot_role_id`,`user_roles`.`user_id` as `goelo_pivot_user_id` from `users` inner join user_roles on `user_roles`.`user_id` = `users`.`id` where `role_id` = ? [10] {2} 60.804572ms}
```
