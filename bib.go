package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
func getBibs(c *gin.Context) {
	formatter, _ := c.GetQuery("f")

	for i := 0; i < len(bibs); i++ {
		bibs[i].Format(formatter, i+1)
	}

	c.IndentedJSON(http.StatusOK, bibs)
}

// POST /bibs 创建新的参考文献，添加到bibs
func postBib(c *gin.Context) {
	var newBib Bibliography

	err := c.Bind(&newBib)
	if err != nil {
		return
	}

	bibs = append(bibs, newBib)

	c.IndentedJSON(http.StatusCreated, newBib)

}

// GET /bibs/:id 获取指定 id 的一篇参考文献，若指定 f=gbt 则格式化为 GB/T 7714 格式
func getBibByID(c *gin.Context) {
	id := c.Param("id")
	formatter, _ := c.GetQuery("f")

	for i, b := range bibs {
		if b.ID == id {
			b.Format(formatter, i+1)
			c.IndentedJSON(http.StatusOK, b)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bib not found"})
}

func main() {
	router := gin.Default()

	router.GET("/bibs", getBibs)
	router.GET("/bibs/:id", getBibByID)
	router.POST("/bibs", postBib)

	router.Run("localhost:8080")
}
