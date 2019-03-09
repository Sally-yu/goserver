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
		db := database.DbConnection{Dbname, Cname, nil, nil, nil}
		img := model.Img{jsono["deviceid"], jsono["imgurl"]}
		err := img.Save(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data := Response{}
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
		img := model.Img{device["deviceid"], ""}
		err, _ := img.Find(database.DbConnection{Dbname, Cname, nil, nil, nil})
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data := Response{}
		data.Status = img.Imgurl
		json.NewEncoder(w).Encode(data)
	}
}

func BackImg(w http.ResponseWriter, r *http.Request) {
	CorsHeader(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var post map[string]string
		json.Unmarshal(body, &post) //json解析
		back := model.Back{}
		err, _ := back.Find(database.DbConnection{Dbname, Cname, nil, nil, nil})
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode(back)
	}
}

func WorkSpace(w http.ResponseWriter, r *http.Request) {
	CorsHeader(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var post = struct {
			Opt string
			Workspace model.WorkSpace
		}{}
		worklist:=[]model.WorkSpace{}
		json.Unmarshal(body, &post) //json解
		var err error
		switch post.Opt {
		case "save":
			err = post.Workspace.Save(database.DbConnection{"docdb", "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(post.Workspace) //response一个workspace
			break
		case "find":
			err, _= post.Workspace.Find(database.DbConnection{"docdb", "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(post.Workspace) //response一个workspace
			break
		case "all":
			err,worklist= post.Workspace.FindAll(database.DbConnection{"docdb", "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(worklist) //response一个workspace
			break
		default:
			break
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
	}
}

func main() {
	http.HandleFunc("/assets/img", PathServer)
	http.HandleFunc("/assets/img/save", SaveImg)
	http.HandleFunc("/assets/img/deviceid", FindImg)
	http.HandleFunc("/assets/img/back", BackImg)
	http.HandleFunc("/workspace", WorkSpace)

	fs := http.FileServer(http.Dir(path))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", fs)) //设备图片
	err := http.ListenAndServe(":8090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Printf("sever start")
}
