---
title: One To Many Morph
---

We already have post and video model, let's add comment feature.  
Post comments,video comments, basiclly are comments. So we just need a comment table with an addtional column
commentable_type to indicate which type of this comment

### Table Structure

```
posts
    id - integer
    title - string
video
    id - integer
    title - string
comment
    id - integer
    content - string
    commentable_id - integer
    commentable_type - string
   
```

### Model Structure

```
type Post struct {
	goeloquent.EloquentModel
	Id        int64     `goelo:"column:id;primaryKey"`
	UserId    int64     `goelo:"column:user_id"`
	Title     string    `goelo:"column:title"`
	Content   string    `goelo:"column:content"`
	User      *User     `goelo:"BelongsTo:UserRelation"`
	Thumbnail Image     `goelo:"MorphOne:ThumbnailRelation"`
	Comments  []Comment `goelo:"MorphMany:CommentsRelation"`
}

func (p *Post) CommentsRelation() *goeloquent.RelationBuilder {
	return p.MorphMany(p, &Comment{}, "commentable_type", "commentable_id", "id")
}
func (p *Post) TableName() string {
	return "posts"
}

type Video struct {
	goeloquent.EloquentModel
	Id          int64     `goelo:"column:id;primaryKey"`
	UserId      int64     `goelo:"column:user_id"`
	Title       string    `goelo:"column:title"`
	Durition    int       `goelo:"column:durition"`
	UploadAt    time.Time `goelo:"column:upload_at"`
	ViewCount   int       `goelo:"column:view_count"`
	PublishedAt time.Time `goelo:"column:published_at"`
	Cover       Image     `goelo:"MorphOne:ImageRelation"`
	Comments    []Comment `goelo:"MorphMany:CommentsRelation"`
	Tags        []Tag     `goelo:"MorphToMany:TagsRelation"`
	User      *User  `goelo:"BelongsTo:UserRelation"`

}
func (v *Video) CommentsRelation() *goeloquent.RelationBuilder {
	return v.MorphMany(v, &Comment{}, "commentable_type", "commentable_id", "id")
}
func (v *Video) TableName() string {
	return "videos"
}

type Comment struct {
	goeloquent.EloquentModel
	Id              int64       `goelo:"column:id;primaryKey"`
	UserId          int64       `goelo:"column:user_id"`
	ParentId        int         `goelo:"column:parent_id"`
	Content         string      `goelo:"column:content"`
	CommentableId   int64       `goelo:"column:commentable_id"`
	CommentableType string      `goelo:"column:commentable_type"`
	UpvoteCount     int         `goelo:"column:upvote_count"`
	DownvoteCount   int         `goelo:"column:downvote_count"`
	CreatedAt       time.Time   `goelo:"column:created_at,timestatmp:create"`
	Commentable     interface{} `goelo:"MorphTo:CommentAbleRelation"`
}

func (c *Comment) CommentAbleRelation() *goeloquent.RelationBuilder {
	return c.MorphTo(c, "commentable_id", "id", "commentable_type")
}
func (c *Comment) TableName() string {
	return "comments"
}
```

MorphMany

1. add an tag ``goelo:"MorphMany:RelationMethod"`` for one to many morph model field type can be slice of struct,pointer
   slice of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `MorphMany`
3. `MorphMany` takes 5 parameter,first one is a pointer of current model,second is a pointer of related model, third is
   related model database field correspond to current(parent) model morph type , 4th is related model database field
   correspond to current(parent) model id, last one is current model field correspond to 4th parameter.

### Usage Example

Use With when retrive

```
var p models.Post
var ps []models.Post
DB.Model(&p).With("Comments").Find(&p,2)
//{select * from `posts` where `id` in (?) limit 1 [2] {1} 61.820344ms}
//{select * from `comments` where `commentable_id` is  not null  and `commentable_type` = ? and `commentable_id` in (?) [post 2] {2} 61.770256ms}
DB.Model(&p).With("Comments").Get(&ps)
//{select * from `posts` [] {6} 62.168008ms}
//{select * from `comments` where `commentable_id` is  not null  and `commentable_type` = ? and `commentable_id` in (?,?,?,?,?,?) [post 2 4 6 8 10 12] {4} 62.450237ms}

```

Use Load Or Directly Call Relation Method

```
p.Load("Comments")
//{select * from `comments` where `commentable_id` is  not null  and `commentable_type` = ? and `commentable_id` in (?) [post 2] {2} 61.67041ms}
var cs []models.Comment
p.CommentsRelation().Get(&cs)
//{select * from `comments` where `commentable_id` = ? and `commentable_id` is  not null  and `commentable_type` = ? [2 post] {2} 70.803116ms}
```
### Morph Many Reverse
```
var c models.Comment
var cs []models.Comment

var p models.Post
var v models.Video
DB.Model(&c).With("Commentable").Find(&c, 2)
//{select * from `comments` where `id` in (?) limit 1 [2] {1} 60.757119ms}
//{select * from `posts` where `id` in (?) [2] {1} 71.229719ms}

DB.Model(&c).With("Commentable").Get(&cs)
//{select * from `comments` [] {7} 70.454251ms}
//{select * from `posts` where `id` in (?,?,?,?) [2 2 4 6] {3} 90.719874ms}
//{select * from `videos` where `id` in (?,?,?) [2 2 6] {1} 73.704877ms}

c.Load("Commentable")
//{select * from `posts` where `id` in (?) [2] {1} 61.203106ms}

if c.CommentableType == "post" {
  c.CommentableRelation().Get(&p)
//{select * from `posts` where `id` = ? [2] {1} 63.048292ms}
} else if c.CommentableType == "video" {
  c.CommentableRelation().Get(&v)
}
```
