package main

import (
	"encoding/json"
	"fmt"
	"goserver/database"
	"goserver/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// 获取大小的接口
type Sizer interface {
	Size() int64
}
type Response struct {
	Status string
}

const (
	path   = "/usr/local/nginx/html/assets/img"
	Dbname = "imgdb"
	Cname  = "imgcollection"
	Idbname = "wsg"
	Icname =  "test"
)

//const path="c:/img"

//跨域头
func CorsHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                                            //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,X-Requested-With,Authorization") //header的类型	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json") //返回数据格式是json
}

//上传保存图片
func PathServer(w http.ResponseWriter, r *http.Request) {
	CorsHeader(w) //优先处理跨域，否则后续函数不会执行
	if "POST" == r.Method {
		fmt.Println(r)
		file, handler, err := r.FormFile("file") //antd上传控件formdata中文件key为file
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}

		filename := handler.Filename //取文件名
		defer file.Close()
		if _, err := os.Stat(path); err == nil {
			fmt.Println("path exists 1", path)
		} else {
			fmt.Println("path not exists ", path)
			err := os.MkdirAll(path, 0711)
			// check again
			if err != nil {
				io.WriteString(w, "Error creating directory")
				return
			}
			if _, err := os.Stat(path); err == nil {
				fmt.Println("path exists 2", path)
			}
		}

		if _, err := os.Stat(path + "/" + filename); err == nil {
			os.Remove(path + "/" + filename) //删除文件
		}
		f, err := os.Create(path + "/" + filename) //创建文件
		if err != nil {
			return
		}
		defer f.Close()
		var data Response
		data.Status = "success"
		json.NewEncoder(w).Encode(data)
		io.Copy(f, file) //复制文件内容
	}
}

//保存到数据库
func SaveImg(w http.ResponseWriter, r *http.Request) {
	CorsHeader(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var jsono map[string]string
		json.Unmarshal(body, &jsono) //json解析
		db:=database.DbConnection{Dbname,Cname,nil,nil,nil}
		img :=model.Img{jsono["deviceid"], jsono["imgurl"]}
		err := img.Save(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data :=Response{}
		data.Status = "success"
		json.NewEncoder(w).Encode(data)
	}
}

//返回imgurl
func FindImg(w http.ResponseWriter, r *http.Request) {
	CorsHeader(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var device map[string]string
		json.Unmarshal(body, &device) //json解析
		img:=model.Img{device["deviceid"],""}
		err,_:=img.Find(database.DbConnection{Dbname,Cname,nil,nil,nil})
		if err!=nil{
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data:=Response{}
		data.Status=img.Imgurl
		json.NewEncoder(w).Encode(data)
	}
}

func BackImg(w http.ResponseWriter, r *http.Request)  {
	CorsHeader(w)
	if "POST"==r.Method{
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var post map[string]string
 		json.Unmarshal(body, &post) //json解析
 		back:=model.Back{}
 		err,_:=back.Find(database.DbConnection{Dbname,Cname,nil,nil,nil})
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode(back)
	}
}
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
							Datastrategy:"InfluxDB",
							Savestrategy:dataserver["savestrategy"],
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
func main() {
	http.HandleFunc("/assets/img", PathServer)
	http.HandleFunc("/assets/img/save", SaveImg)
	http.HandleFunc("/assets/img/deviceid", FindImg)
	http.HandleFunc("/assets/img/back",BackImg)
	http.HandleFunc("/assets/influx/get",GetInfluxlist)
	http.HandleFunc("/assets/influx/edit",UpdateInfluxlist)
	http.HandleFunc("/assets/influx/delete",DelInfluxlist)
	http.HandleFunc("/assets/influx/insert",InsertInfluxlist)
	http.HandleFunc("/assets/opcua/get",GetOpcualist)
	http.HandleFunc("/assets/opcua/insert",InsertOpcualist)

	fs := http.FileServer(http.Dir(path))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", fs)) //设备图片
	err := http.ListenAndServe(":8090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Printf("sever start")
}
