package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	if err != nil {
		return
	}
}

//func main() {
//	//http.HandleFunc("/", handler)
//	http.HandleFunc("/view/", viewHandler)
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}
