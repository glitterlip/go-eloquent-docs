package main

import (
	"fmt"
	goeloquent "github.com/glitterlip/go-eloquent"
	_ "github.com/go-sql-driver/mysql"
	"go-eloquent-doc/models"
	_ "reflect"
	_ "runtime"
	_ "unicode/utf8"
)

var DB *goeloquent.DB

func init() {
	config := map[string]goeloquent.DBConfig{
		"default": { //conncetion with name default is required
			Host:     "127.0.0.1",
			Port:     "3506",
			Database: "eloquent",
			Username: "root",
			Password: "root",
			Driver:   "mysql",
		},
		"chat": {
			Host:     "127.0.0.1",
			Port:     "3506",
			Database: "chat",
			Username: "root",
			Password: "root",
			Driver:   "mysql",
		},
	}

	//alias for goeloquent.Eloquent
	DB = goeloquent.Open(config)
	//set the logger
	goeloquent.Eloquent.SetLogger(func(log goeloquent.Log) {
		fmt.Println(log)
	})
	//register morph model map
	goeloquent.RegistMorphMap(map[string]interface{}{
		"post":  &models.Post{},
		"video": &models.Video{},
		"tag":   &models.Tag{},
	})
	//register model for performance
	goeloquent.RegisterModels([]interface{}{&models.User{}, &models.Post{}, &models.Info{}, &models.Video{}, &models.Tag{}})
}
func main() {
	Relation()
}
