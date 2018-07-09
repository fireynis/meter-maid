package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"strings"
	"net/url"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Name("home").Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	//Should do err check against parse form.
	data := r.Form
	w.Header().Add("Content-Type", "application/json")
	response := "Got it!"
	command := strings.Fields(data.Get("text"))
	if len(command) == 0 {
		go AlertPeopleAndChannels(data)
		response += " Time stored."
	}
	jsonResp, _ := json.Marshal(JsonResponse{Type:"ephemeral", Text:strings.TrimSpace(response)})
	fmt.Fprint(w, string(jsonResp))
}

func AlertPeopleAndChannels(values url.Values) {
	fmt.Println(values)
}

type JsonResponse struct {
	Type string `json:"response_type"`
	Text string `json:"text"`
}
