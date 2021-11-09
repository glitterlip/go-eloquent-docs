package models

import (
	goeloquent "github.com/glitterlip/go-eloquent"
	"time"
)

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
func (v *Video) CommentsRelation() *goeloquent.RelationBuilder {
	return v.MorphMany(v, &Comment{}, "commentable_type", "commentable_id", "id")
}
func (v *Video) TableName() string {
	return "videos"
}
func (v *Video) UserRelation() *goeloquent.RelationBuilder {
	return v.BelongsTo(v, &User{}, "user_id", "id")
}
func (v *Video) ImageRelation() *goeloquent.RelationBuilder {
	return v.MorphOne(v, &Image{}, "imageable_type", "imageable_id", "id")
}
