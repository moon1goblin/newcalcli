package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	client := http.Client{}
	response, err := client.Get("http://127.0.0.1:6969")
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode == 200 {
		fmt.Println("succes")
	}
}
