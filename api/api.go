package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Create(filename, host, port string) {
	raw, err := ioutil.ReadFile(filename)
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

		r.HandleFunc(fmt.Sprintf("/%s/{id}", p), func(w http.ResponseWriter, r *http.Request) {
			switch v := data[p].(type) {
			case []interface{}:
				for _, m := range v {
					id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
					if int64(m.(map[string]interface{})["id"].(float64)) == id {
						json.NewEncoder(w).Encode(m)
						return
					}
				}
				http.NotFound(w, r)
				return
			default:
				http.NotFound(w, r)
			}
		}).Methods("GET")
	}

	withLogging := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(host+":"+port, withLogging))
}
