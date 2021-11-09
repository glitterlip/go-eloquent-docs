package main

import "go-eloquent-doc/models"

func RelationMorphOne() {
	var p models.Post
	var ps []models.Post

	DB.Model(&p).With("Thumbnail").Find(&p, 2)
	DB.Model(&p).With("Thumbnail").Get(&ps)

	p.Load("Thumbnail")
	var i models.Image
	p.ThumbnailRelation().Get(&i)

}
func RelationMorphOneReverse() {
	var i models.Image
	var is []models.Image
	DB.Model(&i).With("Imageable").Find(&i, 4)
	DB.Model(&i).With("Imageable").Get(&is)

	i.Load("Imageable")

	if i.ImageableType == "post" {
		var t models.Post
		i.ImageableRelation().Get(&t)
	} else if i.ImageableType == "video" {
		var t models.Video
		i.ImageableRelation().Get(&t)
	}

}
