package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	inputChannel  chan map[string]string
	outputChannel chan Output
)

func main() {
	inputChannel = make(chan map[string]string)
	outputChannel = make(chan Output)
	go worker()

	router := http.NewServeMux()
	router.HandleFunc("/assesment", AssesmentHandler)
	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	fmt.Println("Server listening on :8000")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
func AssesmentHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]string

	decoder := json.NewDecoder(r.Body)
	r.Header.Add("Content-Type", "application/json")
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	inputChannel <- req
	json.NewEncoder(w).Encode(<-outputChannel)

}
func output(field map[string]string) {
	outputinput := new(Output)
	outputinput.Event = field["ev"]
	outputinput.EventType = field["et"]
	outputinput.AppID = field["id"]
	outputinput.UserID = field["uid"]
	outputinput.MessageID = field["mid"]
	outputinput.PageTitle = field["t"]
	outputinput.PageURL = field["p"]
	outputinput.BrowserLanguage = field["l"]
	outputinput.ScreenSize = field["cs"]
	outputinput.Attributes = make(map[string]Attribute)
	outputinput.UserTraits = make(map[string]Attribute)
	regex := "^atrk.*"
	regex1 := "^uatrk.*"
	attributes, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}
	userTraits, err := regexp.Compile(regex1)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}
	for key, value := range field {
		if attributes.MatchString(key) {
			str := strings.Split(key, "atrk")
			v := "atrv" + str[1]
			t := "atrt" + str[1]
			var atr Attribute
			atr.Value = field[v]
			atr.Type = field[t]
			outputinput.Attributes[value] = atr
		}
		if userTraits.MatchString(key) {
			str := strings.Split(key, "uatrk")
			v := "uatrv" + str[1]
			t := "uatrt" + str[1]
			var atr Attribute
			atr.Value = field[v]
			atr.Type = field[t]
			outputinput.UserTraits[value] = atr
		}
	}
	outputChannel <- *outputinput
}

func worker() {
	for req := range inputChannel {
		output(req)
	}
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}
type Output struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppID           string               `json:"app_id"`
	UserID          string               `json:"user_id"`
	MessageID       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"`
	UserTraits      map[string]Attribute `json:"traits"`
}
