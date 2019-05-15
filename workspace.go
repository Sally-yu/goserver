package main

import (
	"encoding/json"
	"fmt"
	"goserver/database"
	"goserver/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//上传保存图片
func PathServer(w http.ResponseWriter, r *http.Request,p string) {
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
		if _, err := os.Stat(p); err == nil {
			fmt.Println("path exists 1", p)
		} else {
			fmt.Println("path not exists ", p)
			err := os.MkdirAll(p, 0711)
			// check again
			if err != nil {
				io.WriteString(w, "Error creating directory")
				return
			}
			if _, err := os.Stat(p); err == nil {
				fmt.Println("path exists 2", p)
			}
		}

		if _, err := os.Stat(p + "/" + filename); err == nil {
			os.Remove(p + "/" + filename) //删除文件
		}
		f, err := os.Create(p + "/" + filename) //创建文件
		if err != nil {
			return
		}
		defer f.Close()
		var data Response
		data.Status = "success"
		fmt.Println(filename)
		json.NewEncoder(w).Encode(data)
		io.Copy(f, file) //复制文件内容
	}
}

//保存到数据库
func SaveLink(w http.ResponseWriter, r *http.Request) {
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
			Opt       string
			Workspace model.WorkSpace
		}{}
		worklist := []model.WorkSpace{}
		json.Unmarshal(body, &post) //json解
		var err error
		switch post.Opt {
		case "save":
			err = post.Workspace.Save(database.DbConnection{"docdb", "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(post.Workspace) //response一个workspace
			break
		case "find":
			err, _ = post.Workspace.Find(database.DbConnection{"docdb", "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(post.Workspace) //response一个workspace
			break
		case "all":
			println("all")
			err, worklist = post.Workspace.FindAll(database.DbConnection{"docdb", "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(worklist) //response一个workspace
			break
		case "delete":
			err = post.Workspace.Remove(database.DbConnection{"docdb", "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode("delete success") //response一个workspace
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

func Upload(w http.ResponseWriter, r *http.Request) {
	PathServer(w,r,uploadPath)
}

func SaveSvg(w http.ResponseWriter, r *http.Request){
	PathServer(w,r,path)
}

//更新指定的自定义分组
func UpdateCus(w http.ResponseWriter,r *http.Request){
	CorsHeader(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		cus := model.Cus{}
		json.Unmarshal(body, &cus) //json解析
		c:=model.Cus{Divid:cus.Divid}
		err,_:=c.Find(database.DbConnection{"docdb", "cus", nil, nil, nil})
		fmt.Println("c:",c.Svg)
		fmt.Println("cus:",cus.Svg)
		if err==nil{
			a:=append(cus.Svg,c.Svg...) //合并数组，追加svg
			cus.Svg=a
		}
		cus.RemoveRepeat() //svg去重
		for i:=0 ;i<len(cus.Svg);i++{
			cus.Svg[i].Svg=strings.Replace(cus.Svg[i].Svg,".svg","",1)//去.svg后缀
		}
		err= cus.Update(database.DbConnection{"docdb", "cus", nil, nil, nil})
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data := Response{}
		data.Status = "update"
		json.NewEncoder(w).Encode(data)
	}
}

//获取所有的自定义分组
func CusSvg(w http.ResponseWriter,r *http.Request){
	CorsHeader(w)
	if "GET" == r.Method {
		svgList := []model.Cus{}
		svg:=model.Cus{}
		var err error
		err, svgList = svg.FindAll(database.DbConnection{"docdb", "cus", nil, nil, nil})
		json.NewEncoder(w).Encode(svgList) //response一个workspace
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
	}
}

func FindName(w http.ResponseWriter,r *http.Request){
	CorsHeader(w)
	if "POST"==r.Method{
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		work := model.WorkSpace{}
		var name=""
		json.Unmarshal(body, &name) //json解析
		res:=Response{}
		res.Status=work.FindName(database.DbConnection{"docdb","workspace",nil,nil,nil},name)
		json.NewEncoder(w).Encode(res) //response一个workspace
	}
}