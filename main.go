package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Name("home").Methods("POST")

	log.Fatal(http.ListenAndServe(":4390", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.Form
	log.Printf("%v", data)
	w.Header().Add("Content-Type", "application/json")
	jsonResp, _ := json.Marshal(struct {
		Type string `json:"response_type"`
		Text string `json:"text"`
	}{
		Type: "ephemeral",
		Text: fmt.Sprint("Got it!"),
	})
	fmt.Fprint(w, string(jsonResp))
}
