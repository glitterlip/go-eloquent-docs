---
title: Relations
---
Let me bring you the most part I love about Larave.

When define relation method,parameter with prefix `parent` means current model database field , `related` means related
model database field.





### Many To Many Relationships

Table Structure

```
users
    id - integer
    name - string

roles
    id - integer
    name - string

role_user
    user_id - integer
    role_id - integer
```

Model Structure

```
type RoleUser struct {
	goeloquent.EloquentModel
	Id       int64  `goelo:"column:id;primaryKey"`
	UserName string `goelo:"column:username"`
}
type Role struct {
	goeloquent.EloquentModel
	Id          int64       `goelo:"column:id;primaryKey"`
	Name        string      `goelo:"column:name"`
	DisplayName string      `goelo:"column:display_name"`
	Users       []RoleUser  `goelo:"BelongsToMany:UsersRelation"`
}

func (r *Role) UsersRelation() *goeloquent.RelationBuilder {
	return r.BelongsToMany(r, &RoleUser{}, "user_roles", "role_id", "user_id", "id", "id")
}
func (r *Role) TableName() string {
	return "roles"
}

```
1. add an tag ``goelo:"BelongsToMany:RelationMethod"`` for many to one many field type can be struct, pointer of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `BelongsTo`

`BelongsTo` takes 4 parameter,first one is a pointer of current model,second is a pointer of related model, third is
parent(current) model database field correspond to related model,last one is related model field correspond to third
parameter.
#### Usage Exampl

```go

```
