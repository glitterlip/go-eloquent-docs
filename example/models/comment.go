package models

import (
	goeloquent "github.com/glitterlip/go-eloquent"
	"time"
)

type Comment struct {
	goeloquent.EloquentModel
	Id              int64       `goelo:"column:id;primaryKey"`
	UserId          int64       `goelo:"column:user_id"`
	ParentId        int         `goelo:"column:parent_id"`
	Content         string      `goelo:"column:content"`
	CommentableId   int64       `goelo:"column:commentable_id"`
	CommentableType string      `goelo:"column:commentable_type"`
	UpvoteCount     int         `goelo:"column:upvote_count"`
	DownvoteCount   int         `goelo:"column:downvote_count"`
	CreatedAt       time.Time   `goelo:"column:created_at,timestatmp:create"`
	Commentable     interface{} `goelo:"MorphTo:CommentableRelation"`
}

func (c *Comment) CommentableRelation() *goeloquent.RelationBuilder {
	return c.MorphTo(c, "commentable_id", "id", "commentable_type")
}
func (c *Comment) TableName() string {
	return "comments"
}
