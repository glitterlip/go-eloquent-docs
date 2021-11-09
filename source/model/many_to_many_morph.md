---
title: Many To Many Morph
---

## Many To Many Morph

So let's say we have a tag system.A video or post can have many tags,a tag can belong to many videos/posts. That's a
typical many to many morph relation. In order to achieve that, we need a pivot table.Pivot table tagable_id bounds to
video/post.id,pivot table tagable_type indicate which type (post/video), pivot table tag_id bounds to tag.id.

## Table Structure

```
post
    id - integer
    title - string
video 
    id - integer
    title - string
tag
    id - integer
    title - string
tagable
    id - integer
    tag_id - integer
    tagable_id - integer
    tagable_type - string
```

## Model Structure
```go
type Tag struct {
	goeloquent.EloquentModel
	Id     int64   `goelo:"column:id;primaryKey"`
	Name   string  `goelo:"column:name"`
	Count  int     `goelo:"column:count"`
	Videos []Video `goelo:"MorphByMany:VideosRelation"`
	Posts  []Post  `goelo:"MorphByMany:PostsRelation"`
}

func (t *Tag) PostsRelation() *goeloquent.RelationBuilder {
	return t.MorphByMany(t, &Post{}, "tagables", "tag_id", "tagable_id", "id", "id", "tagable_type")
}
func (t *Tag) VideosRelation() *goeloquent.RelationBuilder {
	return t.MorphByMany(t, &Video{}, "tagables", "tag_id", "tagable_id", "id", "id", "tagable_type")
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
	User        *User     `goelo:"BelongsTo:UserRelation"`
}

func (v *Video) TagsRelation() *goeloquent.RelationBuilder {
	return v.MorphToMany(v, &Tag{}, "tagables", "tag_id", "tagable_id", "id", "id", "tagable_type")
}

type Post struct {
	goeloquent.EloquentModel
	Id        int64     `goelo:"column:id;primaryKey"`
	UserId    int64     `goelo:"column:user_id"`
	Title     string    `goelo:"column:title"`
	Content   string    `goelo:"column:content"`
	User      *User     `goelo:"BelongsTo:UserRelation"`
	Thumbnail Image     `goelo:"MorphOne:ThumbnailRelation"`
	Comments  []Comment `goelo:"MorphMany:CommentsRelation"`
	Tags       []Tag      `goelo:"MorphToMany:TagsRelation"`
}
func (p *Post) TagsRelation() *goeloquent.RelationBuilder {
	return p.MorphToMany(p, &Tag{}, "tagables", "tag_id", "tagable_id", "id", "id", "tagable_type")
}

```

`MorphToMany` takes 8 parameter
1. first one is a pointer of current model
2. second is a pointer of related model 
3. third is pivottable
4. 4th is pivottable relatedModelKey
5. 5th is pivotTable currentTable key
6. 6th is current model key
7. 7th is related model key
8. 8th is pivottable relatedType

`MorphByMany` takes 8 parameter
1. first one is a pointer of current model
2. second is a pointer of related model
3. third is pivottable
4. 4th is pivottable currentModelKey
5. 5th is pivotTable relatedModelkey
6. 6th is current model key
7. 7th is related model key
8. 8th is pivottable relatedType

### Usage Example
Use With when retrive

```go

var p models.Post
var vs []models.Video
DB.Model(&p).With("Tags").Find(&p, 2)
//{select * from `posts` where `id` in (?) limit 1 [2] {1} 69.056284ms}
//{select `tags`.*,`tagables`.`tag_id` as `goelo_pivot_tag_id`,`tagables`.`tagable_id` as `goelo_pivot_tagable_id` from `tags` inner join tagables on `tagables`.`tag_id` = `tags`.`id` where `tagables`.`tagable_id` in (?) and `tagable_type` = ? [2 post] {3} 71.667502ms}

DB.Model(&models.Video{}).With("Tags").Get(&vs)
//{select * from `videos` [] {2} 62.957382ms}
//{select `tags`.*,`tagables`.`tag_id` as `goelo_pivot_tag_id`,`tagables`.`tagable_id` as `goelo_pivot_tagable_id` from `tags` inner join tagables on `tagables`.`tag_id` = `tags`.`id` where `tagables`.`tagable_id` in (?,?) and `tagable_type` = ? [2 4 video] {0} 60.784166ms}

p.Load("Tags")
//{select `tags`.*,`tagables`.`tag_id` as `goelo_pivot_tag_id`,`tagables`.`tagable_id` as `goelo_pivot_tagable_id` from `tags` inner join tagables on `tagables`.`tag_id` = `tags`.`id` where `tagables`.`tagable_id` in (?) and `tagable_type` = ? [2 post] {3} 65.90312ms}

var ts []models.Tag
p.TagsRelation().Get(&ts)
//{select `tags`.*,`tagables`.`tag_id` as `goelo_pivot_tag_id`,`tagables`.`tagable_id` as `goelo_pivot_tagable_id` from `tags` inner join tagables on `tagables`.`tag_id` = `tags`.`id` where `tagable_id` = ? and `tagable_type` = ? [2 post] {3} 65.028647ms}
```
