package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MONGODB_URL="127.0.0.1:27017"
)

type Img struct {
	Deviceid string `bson:"deviceid"`
	Imgurl string `bson:"imgurl"`
}

func ConnDB(dbname string,cname string) (*mgo.Database,*mgo.Collection,error)  {
	session,err:=mgo.Dial(MONGODB_URL)
	if err !=nil {
		return nil,nil,err
	}
	//defer session.Close()

	session.SetMode(mgo.Monotonic,true)
	db:=session.DB(dbname)
	collection:=db.C(cname)

	return db,collection,nil
}

func Insert(dbname string,cname string,deviceid string,imgurl string) error{
	_,collection,err:=ConnDB(dbname,cname)
	err=collection.Insert(&Img{deviceid,imgurl})
	if err!=nil {
		return err
	}
	return nil
}

func Update(dbname string,cname string,deviceid string,imgurl string) error{
	_,collection,err:=ConnDB(dbname,cname)
	//img:=Img{}
	err = collection.Update(bson.M{"deviceid": deviceid}, bson.M{"$set": bson.M{"imgurl": imgurl}})
	if err!=nil{
		return err
	}
	return nil
}

func Find(dbname string,cname string,deviceid string)  (error,string){
	result:=Img{}
	_,collection,err:=ConnDB(dbname,cname)
	err = collection.Find(bson.M{"deviceid": deviceid}).One(&result)
	if err!=nil{
		return err,""
	}
	return nil,result.Imgurl
}

