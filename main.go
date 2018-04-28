package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	raw, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	for k, _ := range data {
		p := k // why doing this? see https://stackoverflow.com/a/44045012/4794989
		r.HandleFunc(fmt.Sprintf("/%s", p), func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]interface{}{p: data[p]})
		}).Methods("GET")
	}

	withLogging := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":6000", withLogging))
}
