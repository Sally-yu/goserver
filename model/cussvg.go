package model

import (
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
)

type Svg struct {
	Svg      string `json:"svg" bson:"svg"`
	Deviceid string `json:"deviceid" bson:"deviceid"`
	Status   string `json:"status" bson:"status"`
}

type Cus struct {
	Divid   string `json:"divid" bson:"divid"`
	Name    string `json:"name" bson:"name"`
	Display bool   `json:"display" bson:"display"`
	Svg     []Svg  `json:"svg" bson:"svg"`
}

//返回所有自定义分组信息
func (cus *Cus) FindAll(db database.DbConnection) (error, []Cus) {
	db.ConnDB()
	res := []Cus{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		print(err.Error())
		return err, nil
	}
	return nil, res
}

func (cus *Cus) Save(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&cus)
	if err != nil {
		return err
	}
	return nil
}

func (cus *Cus) Remove(db database.DbConnection) error {
	db.ConnDB()
	db.Collection.Remove(bson.M{"divid": cus.Divid})
	return nil
}

func (cus *Cus) Update(db database.DbConnection) error {
	err := cus.Remove(db)
	if err != nil {
		return err
	}
	err = cus.Save(db)
	if err != nil {
		return err
	}
	return nil
}
