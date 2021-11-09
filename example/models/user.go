package models

import (
	"database/sql"
	goeloquent "github.com/glitterlip/go-eloquent"
)

type User struct {
	goeloquent.EloquentModel
	Id       int64          `goelo:"column:id;primaryKey"`
	UserName sql.NullString `goelo:"column:username"`
	Age      int            `goelo:"column:age"`
	Balance  int            `goelo:"column:balance"`
	Email    sql.NullString `goelo:"column:email"`
	Info     Info           `goelo:"HasOne:InfoRelation"`
	Posts    []Post         `goelo:"HasMany:PostsRelation"`
	Videos   []Video        `goelo:"HasMany:VideosRelation"`
	Roles    []Role         `goelo:"BelongsToMany:RolesRelation"`
}

func (u *User) TableName() string {
	return "users"
}
func (u *User) PostsRelation() *goeloquent.RelationBuilder {
	return u.HasMany(u, &Post{}, "user_id", "id")
}
func (u *User) VideosRelation() *goeloquent.RelationBuilder {
	return u.HasMany(u, &Video{}, "user_id", "id")
}
func (u *User) InfoRelation() *goeloquent.RelationBuilder {
	return u.HasOne(u, &Info{}, "user_id", "id")
}
func (u *User) RolesRelation() *goeloquent.RelationBuilder {
	return u.BelongsToMany(u, &Role{}, "user_roles", "user_id", "role_id", "id", "id")
}
