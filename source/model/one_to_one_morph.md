---
title: one_to_one_morph
---
## One To One 

If you have  a Post model,a Video model. Post model has a Cover Image,Video model has a Thumbnail Image.

In this case , you can create a post_images and video_images. Or you can use a image table with an additional column indicate which type of the image is.This called Morpy One.


### Table Structure
```
posts
    id - integer
    title - string
videos 
    id - integer
    title - string
image
    id - integer
    path - string
    imageable_id  - integer  
    imageable_type - string 
    
```


## Set Up
Before use morph, you need to register a morphmap, you can call `goeloquent.RegistMorphMap` inside init func
```
//register morph model map
goeloquent.RegistMorphMap(map[string]interface{}{
    "post":  &models.Post{},
    "video": &models.Video{},
    "tag":   &models.Tag{},
})
```

### Model Structure

```
type Post struct {
	goeloquent.EloquentModel
	Id        int64  `goelo:"column:id;primaryKey"`
	UserId    int64  `goelo:"column:user_id"`
	Title     string `goelo:"column:title"`
	Content   string `goelo:"column:content"`
	User      *User  `goelo:"BelongsTo:UserRelation"`
	Thumbnail Image  `goelo:"MorphOne:ThumbnailRelation"`
}

func (p *Post) ThumbnailRelation() *goeloquent.RelationBuilder {
    //we will get you image where image.imageable_type="post" and image.imageable_id=post.id
	return p.MorphOne(p, &Image{}, "imageable_type", "imageable_id", "id")
}
func (p *Post) TableName() string {
	return "posts"
}
func (p *Post) UserRelation() *goeloquent.RelationBuilder {
	return p.BelongsTo(p, &User{}, "user_id", "id")
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
    User        *User     `goelo:"BelongsTo:UserRelation"`

}
func (v *Video) TableName() string {
	return "videos"
}

func (v *Video) ImageRelation() *goeloquent.RelationBuilder {
	return v.MorphOne(v, &Image{}, "imageable_type", "imageable_id", "id")
}
func (v *Video) UserRelation() *goeloquent.RelationBuilder {
	return v.BelongsTo(v, &User{}, "user_id", "id")
}

type Image struct {
	goeloquent.EloquentModel
	Id            int64       `goelo:"column:id;primaryKey"`
	Path          string      `goelo:"column:path"`
	ImageableId   int64       `goelo:"column:imageable_id"`
	ImageableType string      `goelo:"column:imageable_type"`
	Imageable     interface{} `goelo:"MorphTo:ImageableRelation"`
}

func (i *Image) ImageableRelation() *goeloquent.RelationBuilder {
	return i.MorphTo(i, "imageable_id", "id", "imageable_type")
}
func (i *Image) TableName() string {
	return "images"
}

```
MorphOne
1. add an tag ``goelo:"MorphOne:RelationMethod"`` for one to one morph model field type can be struct,pointer of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `MorphOne`
3. `MorphOne` takes 5 parameter,first one is a pointer of current model,second is a pointer of related model, third is
   related model database field correspond to current(parent) model morph type , 4th is related model database field correspond to current(parent) model id, last one is current model field correspond to 4th parameter.

MorphTo
1. add an tag ``goelo:"MorphTo:RelationMethod"`` for one to one morph model field type can be struct,pointer of struct
2. define a method that returns a `*goeloquent.RelationBuilder`,inside method,call `MorphTo`
3. `MorphTo` takes 4 parameter,first one is a pointer of current model,second is a pointer of related model, third is
   current(parent) model database field correspond to related model key , 4th is related model type in MorphMap we registered.When we get imageable_type "post" in databse ,we will find which "post" paired model in MorphMap

### Usage Example
Use With when retrive

```
var p models.Post
var ps []models.Post

DB.Model(&p).With("Thumbnail").Find(&p, 2)
//{select * from `posts` where `id` in (?) limit 1 [2] {1} 64.390881ms}
//{select * from `images` where `imageable_id` is  not null  and `imageable_type` = ? and `imageable_id` in (?) [post 2] {1} 66.56607ms}

DB.Model(&p).With("Thumbnail").Get(&ps)
//{select * from `posts` [] {6} 60.657653ms}
//{select * from `images` where `imageable_id` is  not null  and `imageable_type` = ? and `imageable_id` in (?,?,?,?,?,?) [post 2 4 6 8 10 12] {2} 61.786686ms}


```
Use Load Or Directly Call Relation Method

```
p.Load("Thumbnail")
//{select * from `images` where `imageable_id` is  not null  and `imageable_type` = ? and `imageable_id` in (?) [post 2] {1} 61.137089ms}
var i models.Image
p.ThumbnailRelation().Get(&i)
//{select * from `images` where `imageable_id` = ? and `imageable_id` is  not null  and `imageable_type` = ? [2 post] {1} 60.931456ms}
```

## Reverse/MorphTo

### Usage Example

```
var i models.Image
var is []models.Image
DB.Model(&i).With("Imageable").Find(&i,4)
//{select * from `images` where `id` in (?) limit 1 [4] {1} 60.239766ms}
//{select * from `posts` where `id` in (?) [2] {1} 59.779593ms}
DB.Model(&i).With("Imageable").Get(&is)
//{select * from `images` [] {3} 59.796413ms}
//{select * from `videos` where `id` in (?) [2] {1} 62.999856ms}
//{select * from `posts` where `id` in (?,?) [2 4] {2} 59.979136ms}

```
Use Load Or Directly Call Relation Method
```
i.Load("Imageable")
{select * from `posts` where `id` in (?) [2] {1} 60.900912ms}
var t interface{}

if i.ImageableType == "post" {
    var t models.Post
    i.ImageableRelation().Get(&t)
} else if i.ImageableType == "video" {
    var t models.Video
    i.ImageableRelation().Get(&t)
}
//{select * from `posts` where `id` = ? [2] {1} 63.613889ms}

```
