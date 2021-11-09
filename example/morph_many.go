package main

import (
	"fmt"
	goeloquent "github.com/glitterlip/go-eloquent"
	"go-eloquent-doc/models"
)

func RelationMorphMany() {
	var p models.Post
	var ps []models.Post
	DB.Model(&p).With("Comments").Find(&p, 2)
	DB.Model(&p).With("Comments").Get(&ps)

	p.Load("Comments")
	var cs []models.Comment
	p.CommentsRelation().Get(&cs)

}
func RelationMorphManyReverse() {
	var c models.Comment
	var cs []models.Comment

	var p models.Post
	var v models.Video
	DB.Model(&c).With("Commentable").Find(&c, 2)
	DB.Model(&c).With("Commentable").Get(&cs)

	c.Load("Commentable")
	if c.CommentableType == "post" {
		c.CommentableRelation().Get(&p)
	} else if c.CommentableType == "video" {
		c.CommentableRelation().Get(&v)
	}
}

func RelationMorphToMany() {
	var p models.Post
	var vs []models.Video
	DB.Model(&p).With("Tags").Find(&p, 2)
	DB.Model(&models.Video{}).With("Tags").Get(&vs)

	p.Load("Tags")
	var ts []models.Tag
	p.TagsRelation().Get(&ts)
}

func RelationQuery() {
	var u models.User
	DB.Model(&u).With(map[string]func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder{
		"Posts": func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder {
			builder.Where("title", "like", fmt.Sprintf("%%%s%%", "struct")).OrWhere("content", "like", fmt.Sprintf("%%%s%%", "struct"))
			return builder
		},
	}).Find(&u, 4)
	var posts []models.Post
	u.PostsRelation().Where("title", "like", fmt.Sprintf("%%%s%%", "struct")).Get(&posts)

	var us []models.User
	DB.Model(&u).With("Videos", "Posts").Find(&us, []interface{}{4, 6, 8})

	DB.Model(&u).With(map[string]func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder{
		"Posts.Comments": func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder {
			builder.Where("upvote_count", ">", 0)
			return builder
		},
		"Videos.Comments": func(builder *goeloquent.RelationBuilder) *goeloquent.RelationBuilder {
			builder.Where("upvote_count", ">", 0)
			return builder
		},
	}).Find(&u, 4)
}

func RelationPivot() {
	//var r models.Role
	//var rs []models.Role
	//_, _ = goeloquent.Eloquent.Model(&models.Role{}).Where("id", 12).With("Users").WithPivot("status", "id").First(&r)
	//_, _ = goeloquent.Eloquent.Model(&models.Role{}).With("Users").WherePivot("granted_by", 2).WithPivot("granted_by").Get(&rs)
	//fmt.Println(rs[0].Users[0].Pivot["goelo_pivot_granted_by"])
	var us []models.User

	p := &goeloquent.Paginator{
		Items:       &us,
		PerPage:     2,
		CurrentPage: 2,
	}
	q := DB.Model(&models.User{})
	DB.Model(&models.User{}).When(false, func(builder *goeloquent.Builder) {
		q.Where("id",10)
	}).Get(&us)
	q.Paginate(p)
	fmt.Println(p)
	goeloquent.Eloquent.Raw("default")
}
