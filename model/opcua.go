package model

import (
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
)


type Opcua struct {
	Key            string   `json:"key" bson:"key"`
	Datastrategy   string   `json:"datastrategy" bson:"datastrategy"`
	Savestrategy   string   `json:"savestrategy" bson:"savestrategy"`
	Login          string   `json:"login" bson:"login"`
	Servergroup    string   `json:"servergroup" bson:"servergroup"`

}

func (opcua *Opcua) Insert(db database.DbConnection) error {
	db.ConnDB()
	err:=db.Collection.Remove(bson.M{"datastrategy":"InfluxDB"})

	err = db.Collection.Insert(&opcua)
	if err != nil {
		return err
	}
	return nil
}

func (opcua *Opcua) Find(db database.DbConnection) (error, *Opcua) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"datastrategy":"InfluxDB"}).One(&opcua)
	if err != nil {
		return err, nil
	}
	return nil, opcua
}

func (opcua *Opcua) Remove(db database.DbConnection)  error {
	db.ConnDB()
	err:=db.Collection.Remove(bson.M{"key":opcua.Key})
	if err != nil {
		return err
	}
	return nil
}

func (opcua *Opcua) Update(db database.DbConnection) error {
	db.ConnDB()
	err:=db.Collection.Update(bson.M{"key":opcua.Key},opcua)
	if err != nil {
		return err
	}
	return nil
}
