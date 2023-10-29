package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

// Page 定义页面文件结构体
type Page struct {
	Title string
	Body  []byte
}

// 加缓存，我也不知道是不是真的缓存
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// 验证 url 的正则，不会写抄的
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid page url")
	}
	return m[2], nil // title 在下标2
}

// 保存页面文件
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

// 读取页面文件
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// 查看页面文件
func viewHandler(w http.ResponseWriter, r *http.Request) {
	// 记录访问
	log.Println("viewHandler:", r.URL.Path)
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		fmt.Printf("%s", err)
		// 如果页面不存在就重定向到编辑页面
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

// 编辑页面文件
func editHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("editHandle:", r.URL.Path)
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}
	renderTemplate(w, "edit", p)
}

// 保存页面文件
func saveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("saveHandle:", r.URL.Path)
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{title, []byte(body)}
	errSave := p.save()
	if errSave != nil {
		// 给客户端返回错误
		http.Error(w, errSave.Error(), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// html模板读取
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
