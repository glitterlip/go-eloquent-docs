package models

import goeloquent "github.com/glitterlip/go-eloquent"

type Role struct {
	goeloquent.EloquentModel
	Id          int64  `goelo:"column:id;primaryKey"`
	Name        string `goelo:"column:name"`
	DisplayName string `goelo:"column:display_name"`
	Users       []User `goelo:"BelongsToMany:UsersRelation"`
}

func (r *Role) UsersRelation() *goeloquent.RelationBuilder {
	return r.BelongsToMany(r, &User{}, "user_roles", "role_id", "user_id", "id", "id")
}
func (r *Role) TableName() string {
	return "roles"
}
