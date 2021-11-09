package main

import "go-eloquent-doc/models"

func RelationManyToMany() {
	var u models.User
	var us []models.User
	DB.Model(&models.User{}).With("Roles").Find(&u,4)
	DB.Model(&models.User{}).With("Roles").Where("id", "<", 10).Get(&us)

	var uu models.User
	var roles []models.Role
	DB.Model(&uu).Find(&uu,6)
	uu.Load("Roles")

	uu.RolesRelation().Get(&roles)
}
func RelationManyToManyReverse() {
	var r models.Role
	var rs []models.Role
	DB.Model(&models.Role{}).With("Users").Find(&r,10)
	DB.Model(&models.Role{}).With("Users").Get(&rs)

	var rr models.Role
	var users []models.User
	DB.Model(&rr).Find(&rr,10)
	rr.Load("Users")

	rr.UsersRelation().Get(&users)
}