package main

import (
	goeloquent "github.com/glitterlip/go-eloquent"
	"go-eloquent-doc/models"
)

func RelationOneToMany(){
	var u models.User
	_, _ = goeloquent.Eloquent.Model(&models.User{}).With("Posts").Find(&u, 4)
	var users []models.User
	_,_ = DB.Model(&models.User{}).Where("id","<",10).With("Posts").Get(&users)

	var uu models.User
	DB.Model(&models.User{}).Find(&uu,4)
	uu.Load("Posts")
	var posts []models.Post
 	uu.PostsRelation().Get(&posts)
}

func RelationOneToManyReverse(){
	var ps []models.Post
	DB.Model(&models.Post{}).With("User").Get(&ps)

	ps[0].Load("User")
	var user models.User
	ps[0].UserRelation().Get(&user)
}