package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// 记录访问
	log.Println("viewHandler:", r.URL.Path)
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		fmt.Printf("%s", err)
		http.NotFound(w, r)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("editHandle:", r.URL.Path)
	title := r.URL.Path[len("/edit/"):]
	p, errLoadPage := loadPage(title)
	if errLoadPage != nil {
		fmt.Printf("errLoadPage: %s", errLoadPage)
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("saveHandle:", r.URL.Path)
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{title, []byte(body)}
	errSave := p.save()
	if errSave != nil {
		fmt.Printf("errSave Error: %s", errSave)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		fmt.Printf("renderTemplate Error: %s", err)
	} else {
		errExecute := t.Execute(w, p)
		if errExecute != nil {
			fmt.Printf("Execute Error: %s", err)
		}
	}
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
