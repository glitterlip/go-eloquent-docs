---
title: Query Relations
---

Since all relations are defined via methods,and they are retrived by separate sql querys.

You can add other constraints before load relation results from database.

### Add constrains use `with` or relationmethod

```
var u models.User
DB.Model(&u).With(map[string]func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder{
    "Posts": func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder {
        builder.Where("title", "like", fmt.Sprintf("%%%s%%", "struct")).OrWhere("content", "like", fmt.Sprintf("%%%s%%", "struct"))
        return builder
    },
}).Find(&u, 4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 66.970403ms}
//{select * from `posts` where `user_id` is  not null  and `user_id` in (?) and `title` like ? or `content` like ? [4 %struct% %struct%] {4} 73.976216ms}

var posts []models.Post
u.PostsRelation().Where("title", "like", fmt.Sprintf("%%%s%%", "struct")).Get(&posts) 
//{select * from `posts` where `user_id` = ? and `user_id` is  not null  and `title` like ? [4 %struct%] {1} 66.896803ms}
```

### Load mutiple relations

```go
var us []models.User
DB.Model(&u).With("Videos", "Posts").Find(&us, []interface{}{4, 6, 8})
//{select * from `users` where `id` in (?,?,?) [4 6 8] {3} 67.432028ms}
//{select * from `videos` where `user_id` is  not null  and `user_id` in (?,?,?) [4 6 8] {1} 73.787433ms}
//{select * from `posts` where `user_id` is  not null  and `user_id` in (?,?,?) [4 6 8] {2} 68.884461ms}
```

### Load nested relations

```go
DB.Model(&u).With(map[string]func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder{
    "Posts.Comments":func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder {
        builder.Where("upvote_count", ">", 0)
        return builder
    },
    "Videos.Comments":func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder {
        builder.Where("upvote_count", ">", 0)
        return builder
    },
}).Find(&u, 4)
//{select * from `users` where `id` in (?) limit 1 [4] {1} 65.780493ms}
//{select * from `posts` where `user_id` is  not null  and `user_id` in (?) [4] {2} 63.999841ms}
//{select * from `comments` where `commentable_id` is  not null  and `commentable_type` = ? and `commentable_id` in (?,?) and `upvote_count` > ? [post 2 12 0] {1} 63.856649ms}
//{select * from `videos` where `user_id` is  not null  and `user_id` in (?) [4] {1} 65.756311ms}
//{select * from `comments` where `commentable_id` is  not null  and `commentable_type` = ? and `commentable_id` in (?) and `upvote_count` > ? [video 2 0] {1} 68.023121ms}
```

### Query PivotTable

In many to many relations , there is a pivot table, you can add constraints or get it . Use `WherePivot` and `WithPivot`
In previous example we have a user-role many to many relation.Below is pivot table structure

```
user_roles
    id - integer,
    user_id - integer
    role_id - integer
    granted_by - integer
    granted_at - timestamp
    status - integer
```
Query Example
```
var r Role
var rs []Role
_, _ = goeloquent.Eloquent.Model(&Role{}).Where("id", 12).With("Users", "UsersP").WithPivot("status","id").First(&r)
//{select * from `roles` where `id` = ? limit 1 [12] {1} 74.899725ms}
//{select `users`.*,`user_roles`.`role_id` as `goelo_pivot_role_id`,`user_roles`.`user_id` as `goelo_pivot_user_id`,`user_roles`.`status` as `goelo_pivot_status`,`user_roles`.`id` as `goelo_pivot_id` from `users` inner join user_roles on `user_roles`.`user_id` = `users`.`id` where `user_roles`.`role_id` in (?) [12] {1} 80.412644ms}

_, _ = goeloquent.Eloquent.Model(&Role{}).With("Users", "UsersP").WherePivot("granted_by", 2).WithPivot("granted_by").Get(&rs)
//{select * from `roles` [] {4} 66.349528ms}
fmt.Println(rs[0].Users[0].Pivot["goelo_pivot_granted_by"])
//{2 true} sql.NullString

```
