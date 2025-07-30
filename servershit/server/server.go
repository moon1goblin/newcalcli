package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/moon1goblin/newcalcli/cal"
)

func main() {
	bday, err := cal.ProcessDates("bday", "11 sep", "")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(*bday); err != nil {
			log.Fatal(err)
		}
	})

	// we listen and we ~dont judge~ serve
	if err := http.ListenAndServe("127.0.0.1:6969", nil); err != nil {
		log.Fatal(err)
	}
}
