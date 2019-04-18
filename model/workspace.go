package model

import (
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
)

type Node struct {
	Svg string `json:"svg" bson:"svg"`
	Key int `json:"key" bson:"key"`
	Loc string `json:"loc" bson:"loc"`
	Deviceid string `json:"deviceid" bson:"deviceid"`   //关联设备的id
	Status string `json:"status" bson:"status"`  //运行状态指示
}

type Link struct {
	From int `json:"from" bson:"from"`
	To int `json:"to" bson:"to"`
	FromPortId string `json:"fromPortId" bson:"fromPortId"` //与前端页面中绑定的字段key重名，保存连线发出口
	ToPortId string `json:"toPortId" bson:"toPortId"`  //与前端页面中绑定的字段key重名，保存连线结束口
	Points []string `json:"points" bson:"points"`
}

type WorkSpace struct {
	Name string `json:"name" bson:"name"`
	Key string `json:"key" bson:"key"`
	Class string `json:"class" bson:"class"`
	NodeDataArray []Node `json:"nodeDataArray" bson:"nodeDataArray"`
	LinkDataArray []Link `json:"linkDataArray" bson:"linkDataArray"`
}

func (workspc *WorkSpace) Save(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&workspc)
	if err != nil {
		return err
	}
	return nil
}

func (workspc *WorkSpace) FindName(db database.DbConnection,name string) string{
	db.ConnDB()
	db.Collection.Find(bson.M{"name": name}).One(&workspc)
	if len(workspc.Name)>0{
		return "0" //重复
	}
	return "1"
}

func (workspc *WorkSpace) Find(db database.DbConnection) (error, *WorkSpace) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"key": workspc.Key}).One(&workspc)
	if err != nil {
		return err, nil
	}
	return nil, workspc
}

func (workspc *WorkSpace) FindAll(db database.DbConnection) (error, []WorkSpace){
	db.ConnDB()
	res:=[]WorkSpace{}
	err:=db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	return nil, res
}

func (workspc *WorkSpace) Remove(db database.DbConnection)  error {
	db.ConnDB()
	err:=db.Collection.Remove(bson.M{"key":workspc.Key})
	if err != nil {
		return err
	}
	return nil
}

func (workspc *WorkSpace) Update(db database.DbConnection) error {
	err:=workspc.Remove(db)
	if err != nil {
		return err
	}
	err=workspc.Save(db)
	if err != nil {
		return err
	}
	return nil
}