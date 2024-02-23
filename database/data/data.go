package data

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Delete() {

}
func database() {
	db, err := sql.Open("mysql", "liss32:w1350755@tcp(localhost:3306)/test")
	checkErr(err)
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	//插入
	stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec("李洪明", "乐安", "2018-8-19")
	checkErr(err)
	last, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(last)

	//改
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)
	res, err = stmt.Exec("liss32", last)
	checkErr(err)
	row, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(row)
	//查找
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	//删除
	/*	stmt, err = db.Prepare("delete from userinfo where uid=?")
		checkErr(err)

		res, err = stmt.Exec(last)
		checkErr(err)
		row, err = res.RowsAffected()
		checkErr(err)
		fmt.Println(row)*/
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
