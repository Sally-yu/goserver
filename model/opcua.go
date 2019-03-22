package model

import (
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
)


type Opcua struct {
	Key            string   `json:"key" bson:"key"`
	Datastrategy   string   `json:"datastrategy" bson:"datastrategy"`
	Opcstate       string   `json:"opcstate" bson:"opcstate"`
	Opctype        string   `json:"opctype" bson:"opctype"`
	Opchost        string   `json:"opchost" bson:"opchost"`
	Serverurl      string   `json:"serverurl" bson:"serverurl"`
	Interval       string   `json:"interval" bson:"interval"`
	Savestrategy   string   `json:"savestrategy" bson:"savestrategy"`
	Influxhost     string   `json:"influxhost" bson:"influxhost"`
	Influxdatabase string   `json:"influxdatabase" bson:"influxdatabase"`
	Servergroup    string   `json:"servergroup" bson:"servergroup"`
	Login          string   `json:"login" bson:"login"`
}

func (opcua *Opcua) Insert(db database.DbConnection) error {
	db.ConnDB()
	err:=db.Collection.Remove(bson.M{"datastrategy":"OPCUACONFIG"})

	err = db.Collection.Insert(&opcua)
	if err != nil {
		return err
	}
	return nil
}

func (opcua *Opcua) Find(db database.DbConnection) (error, *Opcua) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"datastrategy":"OPCUACONFIG"}).One(&opcua)
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
	err:=db.Collection.Update(bson.M{"datastrategy":"OPCUACONFIG"},opcua)
	if err != nil {
		return err
	}
	return nil
}
