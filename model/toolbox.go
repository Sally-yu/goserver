package model

import (
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
)

type Tool struct {
	Key string `json:"key" bson:"key"`
	Img string `json:"img" bson:"img"`
}


func (tool *Tool) Save(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&tool)
	if err != nil {
		return err
	}
	return nil
}

func (tool *Tool) Find(db database.DbConnection) (error, *Tool) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"img": tool.Img}).One(&tool)
	if err != nil {
		return err, nil
	}
	return nil, tool
}

func (tool *Tool) Remove(db database.DbConnection)  error {
	db.ConnDB()
	err:=db.Collection.Remove(bson.M{"img":tool.Img})
	if err != nil {
		return err
	}
	return nil
}

func (tool *Tool) Update(db database.DbConnection) error {
	err:=tool.Remove(db)
	if err != nil {
		return err
	}
	err=tool.Save(db)
	if err != nil {
		return err
	}
	return nil
}