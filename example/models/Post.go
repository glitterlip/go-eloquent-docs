package models

import goeloquent "github.com/glitterlip/go-eloquent"

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
func (p *Post) CommentsRelation() *goeloquent.RelationBuilder {
	return p.MorphMany(p, &Comment{}, "commentable_type", "commentable_id", "id")
}
func (p *Post) ThumbnailRelation() *goeloquent.RelationBuilder {
	return p.MorphOne(p, &Image{}, "imageable_type", "imageable_id", "id")
}
func (p *Post) TableName() string {
	return "posts"
}
func (p *Post) UserRelation() *goeloquent.RelationBuilder {
	return p.BelongsTo(p, &User{}, "user_id", "id")
}
