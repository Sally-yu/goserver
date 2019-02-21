package main

import (
	"encoding/json"
	"fmt"
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
		var img map[string]string
		json.Unmarshal(body, &img) //json解析
		err := Insert(Dbname, Cname, img["deviceid"], img["imgurl"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		var data Response
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
		err,imgurl:=Find(Dbname,Cname,device["deviceid"])
		if err!=nil{
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		var data Response
		data.Status = imgurl
		json.NewEncoder(w).Encode(data)

	}
}

func main() {
	http.HandleFunc("/assets/img", PathServer)
	http.HandleFunc("/assets/img/save", SaveImg)
	http.HandleFunc("/assets/img/deviceid", FindImg)
	fs := http.FileServer(http.Dir(path))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", fs)) //设备图片
	err := http.ListenAndServe(":8090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Printf("sever start")
}
