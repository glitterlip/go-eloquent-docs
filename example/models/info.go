package models

import (
	goeloquent "github.com/glitterlip/go-eloquent"
	"time"
)

type Info struct {
	goeloquent.EloquentModel
	Id       int64     `goelo:"column:id;primaryKey"`
	UserId   int64     `goelo:"column:user_id"`
	Address  string    `goelo:"column:address"`
	Gender   int       `goelo:"column:gender"`
	Phone    string    `goelo:"column:phone"`
	Birthday time.Time `goelo:"column:birthday"`
	User     *User     `goelo:"BelongsTo:UserRelation"`

}
func (i *Info) TableName() string {
	return "infos"
}
func (i *Info) UserRelation() *goeloquent.RelationBuilder {
	return i.BelongsTo(i, &User{}, "user_id", "id")
}