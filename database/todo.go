package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sever/database/data"
	"strconv"

	"golang.org/x/net/websocket"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func pre(f func(w http.ResponseWriter, r *http.Request, db *sql.DB)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", "liss32:w1350755@/test")
		checkErr(err)
		f(w, r, db)
	}
}
func render(filepath string, w http.ResponseWriter) {
	filename := filepath + ".html"
	t, err := template.ParseFiles(filename)
	checkErr(err)
	t.Execute(w, nil)
}
func Show(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	/*	if k, ok := r.Header["Content-Type"]; ok {
				//r.Header["Content-Type"][0] == "application/json"
				fmt.Println(k)
				fmt.Println("json")

				fmt.Println(r.Header["Content-Type"])
		} else {
				t, err := template.ParseFiles("../html/list.html")
				checkErr(err)

				t.Execute(w, nil)
				fmt.Println("sdafasff")
				/*用下面这个方法为什么无法输出HTML文件到客户端？
				t := template.New("some template") //创建一个模板
		   		t, _ = t.ParseFiles("tmpl/welcome.html", nil)  //解析模板文件
		   		user := GetUser() //获取当前用户信息
		    	t.Execute(w, user)  //执行模板的merger操作

				return
			}*/

	/*	row, err := db.Query("SELECT * FROM todo")
		checkErr(err)

		var list data.Todolist
		for row.Next() {
			var iop data.Tododata
			row.Scan(&iop.Id, &iop.Status, &iop.Body)
			list = append(list, iop)
		}
		fmt.Println(list)*/
	or := orm.NewOrm()
	var todo data.Tododata
	var list []*data.Tododata

	qs := or.QueryTable(todo)
	qs.All(&list)

	js, err := json.Marshal(list)
	checkErr(err)

	w.Header().Add("Content-type", "application/json")
	fmt.Println(string(js))
	fmt.Fprint(w, string(js))
}
func Add(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method == "GET" {
		t, err := template.ParseFiles("./html/add.html")
		checkErr(err)

		t.Execute(w, nil)

	} else {
		stmt, err := db.Prepare("INSERT todo set Status=?,Body=?")
		checkErr(err)

		body := r.FormValue("body")

		res, err := stmt.Exec("未完成", body)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)

		fmt.Println(id)
		fmt.Fprintf(w, "已经添加")
	}
}

func predb() *sql.DB {
	db, err := sql.Open("mysql", "liss32:w1350755@/test")
	checkErr(err)
	return db
}
func DeleteTaskFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	orm := orm.NewOrm()
	t := data.Tododata{Id: 4, Status: "未完成"}
	id, err := orm.Insert(&t)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)
	checkErr(err)
	fmt.Printf("insert id: %d", id)

	t.Status = "慢即是快"
	id, err = orm.Update(&t)
	checkErr(err)
	fmt.Printf("update id: %d", id)

	todo := data.Tododata{Id: 4}
	err = orm.Read(&todo)
	fmt.Printf("ERR: %v\n", err)
	fmt.Println(todo)

	id, err = orm.Delete(&todo)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)
	/*
		stmt, err := db.Prepare("DELETE FROM todo where Id=?")
		checkErr(err)

		id := r.FormValue("id")
		rows, err := stmt.Exec(id)
		checkErr(err)

		affect, err := rows.RowsAffected()
		checkErr(err)

		fmt.Println(affect)
	*/
}
func filesever(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/static/"):]
	path = "./html/" + path

	//验证有效性

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		panic(err)
	}

	buf := make([]byte, 512)
	_, err = file.ReadAt(buf, 0)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

	contenttype := http.DetectContentType(buf)
	w.Header().Add("Content-Type", contenttype)

	body, err := io.ReadAll(file)
	checkErr(err)

	_, err = w.Write(body)
	checkErr(err)

}
func EditTaskFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "GET" {
		render("html/edit", w)
	} else {
		var todo data.Tododata
		todo.Body = r.FormValue("body")
		var err error
		todo.Id, err = strconv.Atoi(r.FormValue("id"))
		checkErr(err)
		todo.Status = "未完成"
		affect, err := orm.NewOrm().Update(&todo)
		checkErr(err)
		fmt.Println(affect)
		/*stmt, err := db.Prepare("update todo set Body=? where Id=?")
		checkErr(err)

		res, err := stmt.Exec(body, id)
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)*/

	}
}
func CompleteTaskFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func Query(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func Echo(ws *websocket.Conn) {
	var err error
	fmt.Println("out of for")
	for {
		var reply string
		fmt.Println("in  for")
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println(err.Error())
			fmt.Println("Can't receive")
			fmt.Printf("/n")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func init() {
	orm.RegisterDataBase("default", "mysql", "liss32:w1350755@/test", 30)

	orm.RegisterModel(new(data.Tododata))

	orm.RunSyncdb("default", false, true)
}

func main() {

	http.HandleFunc("/show", pre(Show))
	http.HandleFunc("/add/", pre(Add))
	http.HandleFunc("/delete/", pre(DeleteTaskFunc))
	http.HandleFunc("/edit/", pre(EditTaskFunc))
	http.HandleFunc("/complete/", pre(CompleteTaskFunc))
	http.HandleFunc("/static/", filesever)
	http.Handle("/", websocket.Handler(Echo))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
