package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

type Uploaddir []fs.DirEntry

func loaddir(name string) Serverslice {
	f, err := os.Open("../upload/" + name)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()
	d, error := f.ReadDir(-1)
	if error != nil {
		fmt.Println(error.Error())
	}
	var file Serverslice
	for _, v := range d {
		file.Servers = append(file.Servers, Up{Name: v.Name(), Filetype: v.Name()})
	}
	return file
}

type Up struct {
	Name     string
	Filetype string
}

type Serverslice struct {
	Servers []Up
}

func inputjson(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/backup/"):]
	u := loaddir(title)
	f, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprintf(w, string(f))
}
func show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
func main() {
	http.HandleFunc("/backup", show)
	http.HandleFunc("/backup/", inputjson)

	log.Fatal(http.ListenAndServe(":9090", nil))

}
