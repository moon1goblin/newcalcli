package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
)

func main() {
	client := http.Client{}
	response, err := client.Get("http://127.0.0.1:6969")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("%s\n", body)
}
