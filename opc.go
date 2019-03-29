package main

import (
	"encoding/json"
	"fmt"
	"goserver/database"
	"goserver/model"
	"io/ioutil"
	"net/http"
)

func InsertInfluxlist(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "POST"==r.Method{
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var server map[string] string
		json.Unmarshal(body, &server) //json解析
		influx:=model.InfluxList{Key:model.UniqueId(),
			Servername:server["servername"],
			Serveraddress:server["serveraddress"],
			Database:server["database"],
			Databasetype:server["databasetype"],
			Username:server["username"],
			Password:server["password"]}
		db:=database.DbConnection{Idbname,Icname,nil,nil,nil}
		err:=influx.Insert(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode("")
	}
}
func UpdateInfluxlist(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "POST"==r.Method{
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var server map[string] string
		json.Unmarshal(body, &server) //json解析
		influx:=model.InfluxList{   server["key"],
			server["servername"],
			server["serveraddress"],
			server["database"],
			server["databasetype"],
			server["username"],
			server["password"]}
		db:=database.DbConnection{Idbname,Icname,nil,nil,nil}
		err:=influx.Update(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode("")
	}
}
func DelInfluxlist(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "POST"==r.Method{
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		//var serverip map[bson.ObjectId]bson.ObjectId
		var serverid map[string]string
		json.Unmarshal(body, &serverid) //json解析
		influx:=model.InfluxList{Key:serverid["key"]}
		db:=database.DbConnection{Idbname,Icname,nil,nil,nil}
		err:=influx.Remove(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode("OK")
	}
}
func GetInfluxlist(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "GET"==r.Method{
		influx:=model.InfluxList{}
		db:=database.DbConnection{Idbname,Icname,nil,nil,nil}
		err,result:=influx.Find(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode(result)
	}
}
func GetOpcualist(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "POST"==r.Method{
		opcua:=model.Opcua{ }
		db:=database.DbConnection{Idbname,"server",nil,nil,nil}
		err,opcresult:=opcua.Find(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode(opcresult)
	}
}
func InsertOpcualist(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "POST"==r.Method{
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var dataserver map[string]string
		json.Unmarshal(body, &dataserver) //json解析
		opcua:=model.Opcua{ Key:model.UniqueId(),
			Datastrategy:"OPCUACONFIG",
			Opctype:dataserver["opctype"],
			Opchost:dataserver["opchost"],
			Serverurl:dataserver["serverurl"],
			Opcstate:dataserver["opcstate"],
			Interval:dataserver["interval"],
			Savestrategy:dataserver["savestrategy"],
			Influxhost:dataserver["influxhost"],
			Influxdatabase:dataserver["influxdatabase"],
			Login:dataserver["login"],
			Servergroup:dataserver["servergroup"]}
		db:=database.DbConnection{Idbname,"server",nil,nil,nil}
		err:=opcua.Insert(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode("OK")
	}
}
func UpdateOpcualist(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "POST"==r.Method{
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var dataserver map[string] string
		json.Unmarshal(body, &dataserver) //json解析
		opcua:=model.Opcua{ Key:model.UniqueId(),
			Datastrategy:"OPCUACONFIG",
			Opctype:dataserver["opctype"],
			Opchost:dataserver["opchost"],
			Opcstate:dataserver["opcstate"],
			Serverurl:dataserver["serverurl"],
			Interval:dataserver["interval"],
			Savestrategy:dataserver["savestrategy"],
			Influxhost:dataserver["influxhost"],
			Influxdatabase:dataserver["influxdatabase"],
			Login:dataserver["login"],
			Servergroup:dataserver["servergroup"]}
		db:=database.DbConnection{Idbname,"server",nil,nil,nil}
		err:=opcua.Update(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode("")
	}
}