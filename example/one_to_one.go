package main

import "go-eloquent-doc/models"

func RelationOneToOne() {
	var u models.User
	_, _ = DB.Model(&models.User{}).With("Info").Find(&u, 4)

	var user []models.User
	_, _ = DB.Model(&models.User{}).Where("id", "<", 10).With("Info").Get(&user)

	var e models.User
	DB.Model(&e).Find(&e, 4)
	e.Load("Info")

	var info models.Info
	e.InfoRelation().Get(&info)

}

func RelationOneToOneReverse() {
	var info models.Info
	_, _ = DB.Model(&models.Info{}).With("User").Find(&info, 4)

	var infos []models.Info
	_, _ = DB.Model(&models.Info{}).Where("id", "<", 10).With("User").Get(&infos)

	var i models.Info
	DB.Model(&i).Find(&i, 4)
	i.Load("User")

	var user models.User
	i.UserRelation().Get(&user)

}
