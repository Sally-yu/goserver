package model

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
	"io"
)



type InfluxList struct {
	Key            string   `json:"key" bson:"key"`
	Servername     string   `json:"servername" bson:"servername"`
	Serveraddress  string   `json:"serveraddress" bson:"serveraddress"`
	Database       string	`json:"database" bson:"database"`
	Databasetype   string	`json:"databasetype" bson:"databasetype"`
	Username       string	`json:"username" bson:"username"`
	Password       string	`json:"password" bson:"password"`
}
//type Wsg []InfluxList

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
//生成Guid字串
func  UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func (influx *InfluxList) Insert(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&influx)
	if err != nil {
		return err
	}
	return nil
}

func (influx *InfluxList) Find(db database.DbConnection) (error, []InfluxList) {
	db.ConnDB()
	result:=[]InfluxList{}
	err := db.Collection.Find(nil).All(&result)
	if err != nil {
		return err, nil
	}
	return nil, result
}

func (influx *InfluxList) Remove(db database.DbConnection)  error {
	db.ConnDB()
	err:=db.Collection.Remove(bson.M{"key":influx.Key})
	if err != nil {
		return err
	}
	return nil
}

func (influx *InfluxList) Update(db database.DbConnection) error {
	db.ConnDB()
	err:=db.Collection.Update(bson.M{"key":influx.Key},influx)
	if err != nil {
		return err
	}
	return nil
}
