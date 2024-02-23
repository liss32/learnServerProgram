package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		err := r.ParseMultipartForm(1 << 62)
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
		}
		file, fileheader, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Println(fileheader.Filename)
		fmt.Println(fileheader.Header)
		fmt.Println(r.MultipartForm)

		fmt.Println(fileheader.Size)
		f, err := os.OpenFile("./test/"+fileheader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
		}
		io.Copy(f, file)
		fmt.Fprintf(w, fileheader.Filename+"is saved")
	}

}

func main() {

	http.Handle("/static", http.FileServer(http.Dir("D:\vr")))
	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(":9090", nil))

}
