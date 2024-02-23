package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type page struct {
	Title string
	Body  []byte
}
type account struct {
	username string
	passward string
}

func loadaccount() account {
	message, err := os.ReadFile("./html/account.txt")
	fmt.Println(string(message))
	if err != nil {
		fmt.Println(err.Error())
	}
	u, p, b := strings.Cut(string(message), " ")
	if !b {
		fmt.Println("load account error")
	}
	fmt.Printf("%d %d   ", len(u), len(p))
	fmt.Println(u + "  " + p)
	acc := account{u, p}
	return acc

}

func (p *page) save() error {
	Title := "./html/" + p.Title + ".txt"
	return os.WriteFile(Title, p.Body, 0600)
}
func loadpage(Title string) (*page, error) {
	txt := "./html/" + Title + ".txt"
	Body, err := os.ReadFile(txt)
	if err != nil {
		return nil, err
	}
	return &page{Title: Title, Body: Body}, nil
}

/*
	func (p *page) loadpage(Title string) error {
		txt := Title + ".txt"
		Body, err := os.ReadFile(txt)
		if err != nil {
			return err
		}
		p.Title = Title
		p.Body = Body
		return nil
	}
*/
func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	t, err := template.ParseFiles("./html/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	error := t.Execute(w, p)
	if error != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func viewhandler(w http.ResponseWriter, r *http.Request) {
	LearnURL(r)
	txt := r.URL.Path[len("/view/"):]
	p, err := loadpage(txt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "view", p)
}

func edithandler(w http.ResponseWriter, r *http.Request) {
	LearnURL(r)
	txt := r.URL.Path[len("/edit/"):]
	p, err := loadpage(txt)
	if err != nil {
		p = &page{Title: txt}
	}
	renderTemplate(w, "edit", p)

}
func savehandler(w http.ResponseWriter, r *http.Request) {
	LearnURL(r)
	txt := r.URL.Path[len("/save/"):]
	Body := r.FormValue("Body")
	p := &page{Title: txt, Body: []byte(Body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+p.Title, http.StatusFound)
}
func LearnURL(r *http.Request) {
	r.ParseForm()
	fmt.Println("r.URL.Path:" + r.URL.Path)
	fmt.Println("r.URL.Scheme:" + r.URL.Scheme)
	fmt.Println("r.RequestURI:" + r.RequestURI)
	fmt.Println("r.postForm：", r.PostForm)
	fmt.Println("r.Form：", r.Form)
	fmt.Println("Proto version:" + r.Proto)
	fmt.Println(r.Method)
	fmt.Println()
	for k, v := range r.Form {
		fmt.Println("key:" + k + " value:" + strings.Join(v, " "))
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	if r.Method == "GET" {
		renderTemplate(w, "login", nil)
	} else {
		acc := loadaccount()
		r.ParseForm()
		LearnURL(r)
		if r.Form.Get("username") == acc.username && r.Form.Get("password") == acc.passward {
			http.Redirect(w, r, "view/type", http.StatusFound)
		} else {
			fmt.Fprintf(w, "error username or password")
		}
	}

}

func main() {
	maxFormSize := int64(10 << 20)
	fmt.Println(maxFormSize)
	http.HandleFunc("/login", login)
	http.HandleFunc("/view/", viewhandler)
	http.HandleFunc("/edit/", edithandler)
	http.HandleFunc("/save/", savehandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
