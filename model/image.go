package model

import (
	"gopkg.in/mgo.v2/bson"
	"helloDocker/database"
)

type Img struct {
	Deviceid string `json:"deviceid" bson:"deviceid"`
	Imgurl   string `json:"imgurl" bson:"imgurl"`
}

func (img Img) Save(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&img)
	if err != nil {
		return err
	}
	return nil
}

func (img *Img) Find(db database.DbConnection) (error, *Img) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"deviceid": img.Deviceid}).One(&img)
	if err != nil {
		return err, nil
	}
	return nil, img
}

func (img *Img) Remove(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Remove(bson.M{"deviceid": img.Deviceid})
	if err != nil {
		return err
	}
	return nil
}

func (img *Img) Update(db database.DbConnection) error {
	err := img.Remove(db)
	if err != nil {
		return err
	}
	err = img.Save(db)
	if err != nil {
		return err
	}
	return nil
}
