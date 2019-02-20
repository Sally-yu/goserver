package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 获取大小的接口
type Sizer interface {
	Size() int64
}

const path  = "/usr/local/nginx/html/assets/img"
//const path="c:/img"

// hello world, the web server
func PostServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                                            //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,X-Requested-With,Authorization") //header的类型	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

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
		defer f.Close()
		var data= struct {

		}{}
		json.NewEncoder(w).Encode(data)
		io.Copy(f, file) //复制文件内容
		return
	}
}


func main() {
	http.HandleFunc("/assets/img", PostServer)
	fs := http.FileServer(http.Dir(path))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", fs))

	err := http.ListenAndServe(":8090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Printf("sever start")
}


