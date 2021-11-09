---
title: Pagination
---

We  have a Paginate method

### Usage Example
```
var us []models.User

p := &goeloquent.Paginator{
    Items:       &us,
    PerPage:     2,
    CurrentPage: 2,
}
q := DB.Model(&models.User{})
q.When(false, func(builder *goeloquent.Builder) {
    q.Where("id",10)
})
q.Paginate(p)

//{select count(*) as aggregate from `users` [] {1} 66.865138ms}
//{select * from `users` limit 2 offset 2 [] {2} 63.524782ms}
```
