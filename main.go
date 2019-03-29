package main

import (
	"log"
	"net/http"
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
	uploadPath="/usr/local/nginx/html/assets/upload"
	Dbname = "imgdb"
	Cname  = "imgcollection"
	Idbname = "wsg"
	Icname =  "test"
)

//const path="c:/img"

//跨域头
func CorsHeader(w http.ResponseWriter) http.ResponseWriter{
	println("CorsHeader")
	w.Header().Set("Access-Control-Allow-Origin", "*")                                            //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,X-Requested-With,Authorization,text/html") //header的类型	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json") //返回数据格式是json
	println("Cors Done")
	return w
}

func main() {
	//预制图标存取
	http.HandleFunc("/assets/img", SaveSvg)          //保存svg 自带get文件服务器
	http.HandleFunc("/assets/img/save", SaveLink)    //保存设备和svg的联系
	http.HandleFunc("/assets/upload", Upload)        //上传自定义svg ，自带get文件服务器
	http.HandleFunc("/assets/img/deviceid", FindImg) //加载图片
	http.HandleFunc("/assets/img/back", BackImg)

	//布局、自定义图标存取
	http.HandleFunc("/assets/img/cussvg", CusSvg)    //get上传的自定义svg
	http.HandleFunc("/assets/updateCus", UpdateCus)  //更新自定义信息
	http.HandleFunc("/workspace", WorkSpace)         //保存工作区
	http.HandleFunc("/workspace/findname", FindName) //保存工作区

	//opc、influx
	http.HandleFunc("/assets/influx/get", GetInfluxlist)
	http.HandleFunc("/assets/influx/edit", UpdateInfluxlist)
	http.HandleFunc("/assets/influx/delete", DelInfluxlist)
	http.HandleFunc("/assets/influx/insert", InsertInfluxlist)
	http.HandleFunc("/assets/opcua/get", GetOpcualist)
	http.HandleFunc("/assets/opcua/insert", InsertOpcualist)
	http.HandleFunc("/assets/opcua/update", UpdateOpcualist)

	fs := http.FileServer(http.Dir("/usr/local/nginx/html/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs)) //开启assets文件夹服务
	println("server start…")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}