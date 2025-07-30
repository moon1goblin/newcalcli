package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/moon1goblin/newcalcli/cal"
)

func main() {
	client := http.Client{}
	response, err := client.Get("http://127.0.0.1:6969")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	bday := cal.Event{}
	if err := decoder.Decode(&bday); err != nil {
		log.Fatal(err)
	}

	fmt.Println(bday.String(false))
	fmt.Println(bday.Begin_time)
}
