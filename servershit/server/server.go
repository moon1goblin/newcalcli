package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hey bitch")
	})
	// we listen and we ~dont judge~ serve
	if err := http.ListenAndServe("127.0.0.1:6969", nil); err != nil {
		log.Fatal(err)
	}
}
