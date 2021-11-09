package models

import goeloquent "github.com/glitterlip/go-eloquent"

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
func (t *Tag) TableName() string {
	return "tags"
}
