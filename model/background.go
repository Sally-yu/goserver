package model

import (
	"gopkg.in/mgo.v2/bson"
	"helloDocker/database"
)

type Position struct {
	Top  int
	Left int
}

type Back struct {
	Backname  string  `json:"backname" bson:"backname"`
	Images    [] string  `json:"images" bson:"images"`
	Locations [] Position	`json:"locations" bson:"locations"`
	Width     int	`json:"width" bson:"width"`
	Height    int	`json:"height" bson:"height"`
}

func (back *Back) Save(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&back)
	if err != nil {
		return err
	}
	return nil
}

func (back *Back) Find(db database.DbConnection) (error, *Back) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"backname": back.Backname}).One(&back)
	if err != nil {
		return err, nil
	}
	return nil, back
}

func (back *Back) Remove(db database.DbConnection)  error {
	db.ConnDB()
	err:=db.Collection.Remove(bson.M{"backname":back.Backname})
	if err != nil {
		return err
	}
	return nil
}

func (back *Back) Update(db database.DbConnection) error {
	err:=back.Remove(db)
	if err != nil {
		return err
	}
	err=back.Save(db)
	if err != nil {
		return err
	}
	return nil
}