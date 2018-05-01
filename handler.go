package jhop

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// NewHandler generates new http.Handler, which handles
// resources from given files.
func NewHandler(rs ...io.Reader) (http.Handler, error) {
	router := mux.NewRouter()
	for _, r := range rs {
		var resources map[string]interface{}
		if err := json.NewDecoder(r).Decode(&resources); err != nil {
			return nil, errors.Wrap(err, "unmarshal failed")
		}

		addResource(router, resources)
	}

	return router, nil
}

func addResource(r *mux.Router, resources map[string]interface{}) {
	for k := range resources {
		p := k // why doing this? see https://stackoverflow.com/a/44045012/4794989
		r.HandleFunc(fmt.Sprintf("/%s", p), func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]interface{}{p: resources[p]})
		}).Methods("GET")

		r.HandleFunc(fmt.Sprintf("/%s/{id}", p), func(w http.ResponseWriter, r *http.Request) {
			switch v := resources[p].(type) {
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
}
