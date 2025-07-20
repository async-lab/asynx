package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World! test")
	})

	fmt.Println("Server started: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
