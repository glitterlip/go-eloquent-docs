package models

import goeloquent "github.com/glitterlip/go-eloquent"

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