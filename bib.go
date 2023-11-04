package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Bibliography 参考文献
type Bibliography struct {
	ID        string
	Author    string
	Title     string
	Type      string //文献类型
	Publisher string //出版社
	Reference string //参考文献

	Comment   string //我的注释
	Formatted string //格式化文本
}

// Format 格式化参考文献，将结果填充到formatter里去
func (b *Bibliography) Format(formatter string, idx int) {
	switch formatter {
	case "gbt":
		b.Formatted = fmt.Sprintf("[%v] %v, %v[%v], %v:%v",
			idx, b.Author, b.Title, b.Type, b.Publisher, b.Reference)
	default:
		b.Formatted = ""
	}
}

var bibs = []Bibliography{
	{
		ID:        "1",
		Author:    "John Smith",
		Title:     "Introduction to Computer Science",
		Type:      "Book",
		Publisher: "ABC Publishing",
		Reference: "Smith, J. (2020). Introduction to Computer Science. ABC Publishing.",
		Comment:   "A great introductory book for beginners.",
		Formatted: "Smith, J. (2020). Introduction to Computer Science. ABC Publishing.",
	},
	{
		ID:        "2",
		Author:    "Alice Johnson",
		Title:     "Programming Languages and Their Evolution",
		Type:      "Journal Article",
		Publisher: "Tech Journal",
		Reference: "Johnson, A. (2019). Programming Languages and Their Evolution. Tech Journal.",
		Comment:   "This article provides a comprehensive overview of programming languages.",
		Formatted: "Johnson, A. (2019). Programming Languages and Their Evolution. Tech Journal.",
	},
	{
		ID:        "3",
		Author:    "David Brown",
		Title:     "Machine Learning for Beginners",
		Type:      "E-book",
		Publisher: "Online Publications",
		Reference: "Brown, D. (2021). Machine Learning for Beginners. Online Publications.",
		Comment:   "An excellent resource for those new to machine learning.",
		Formatted: "Brown, D. (2021). Machine Learning for Beginners. Online Publications.",
	},
	{
		ID:        "4",
		Author:    "Emily White",
		Title:     "History of Art",
		Type:      "Textbook",
		Publisher: "Art Press",
		Reference: "White, E. (2018). History of Art. Art Press.",
		Comment:   "A comprehensive textbook for art history enthusiasts.",
		Formatted: "White, E. (2018). History of Art. Art Press.",
	},
	{
		ID:        "5",
		Author:    "Michael Turner",
		Title:     "Research Methods in Social Sciences",
		Type:      "Book",
		Publisher: "Social Science Publishers",
		Reference: "Turner, M. (2017). Research Methods in Social Sciences. Social Science Publishers.",
		Comment:   "A valuable resource for social science researchers.",
		Formatted: "Turner, M. (2017). Research Methods in Social Sciences. Social Science Publishers.",
	},
}

// GET /bibs[?f=gbt] 获取所有参考文献，若指定 f=gbt 则格式化为GB/T 7714 格式
func getBibs(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.String())
	fomatter := r.URL.Query().Get("f")

	for i := 0; i < len(bibs); i++ {
		bibs[i].Format(fomatter, i+1)
	}

	j, _ := json.MarshalIndent(bibs, "", "    ")
	_, err := w.Write(j)
	if err != nil {
		return
	}
}

// POST /bibs 创建新的参考文献，添加到bibs
func postBib(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.String())

	var newBib Bibliography

	body, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(body, &newBib); err != nil {
		return
	}

	bibs = append(bibs, newBib)

	j, _ := json.MarshalIndent(newBib, "", "   ")

	_, err := w.Write(j)
	if err != nil {
		return
	}
}

// GET /bibs/:id 获取指定 id 的一篇参考文献，若指定 f=gbt 则格式化为 GB/T 7714 格式
func getBibByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/bibs/"):]
	formater := r.URL.Query().Get("f")

	log.Println(r.Method, r.URL.String(), fmt.Sprintf("id=%v", id))

	for i, b := range bibs {
		if b.ID == id {
			b.Format(formater, i+1)
			j, _ := json.MarshalIndent(b, "", "  ")
			w.Write(j)
			return
		}
	}

	j, _ := json.MarshalIndent(
		map[string]string{"message": "bib not found"},
		"", "  ")

	_, err := w.Write(j)
	if err != nil {
		return
	}
}

func handleBibs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBibs(w, r)
	case http.MethodPost:
		postBib(w, r)
	}
}

func main() {
	http.HandleFunc("/bibs/", getBibByID)
	http.HandleFunc("/bibs", handleBibs)

	log.Fatal(
		http.ListenAndServe("localhost:8080", nil))
}
