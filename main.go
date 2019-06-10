package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type keptnEvent struct {
	Specversion     string `json:"specversion"`
	Type            string `json:"type"`
	Source          string `json:"source"`
	ID              string `json:"id"`
	Time            string `json:"time"`
	Datacontenttype string `json:"datacontenttype"`
	Shkeptncontext  string `json:"shkeptncontext"`
	Data            struct {
		Githuborg          string `json:"githuborg"`
		Project            string `json:"project"`
		Teststrategy       string `json:"teststrategy"`
		Deploymentstrategy string `json:"deploymentstrategy"`
		Stage              string `json:"stage"`
		Service            string `json:"service"`
		Image              string `json:"image"`
		Tag                string `json:"tag"`
		EvaluationPassed   bool   `json:"evaluationpassed,omitempty"`
	} `json:"data"`
}

var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

var ufoAddress string
var ufoRow string

//Logging : sets up info and error logging
func Logging(infoLogger io.Writer, errorLogger io.Writer) {
	infoLog = log.New(infoLogger, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(errorLogger, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

//keptnHandler : receives keptn events via http and sets UFO LEDs based on payload
func keptnHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event keptnEvent
	err := decoder.Decode(&event)
	if err != nil {
		fmt.Println("Error while parsing JSON payload: " + err.Error())
		return
	}

	if event.Data.Stage == "dev" {
		ufoRow = "top"
	} else if event.Data.Stage == "staging" {
		ufoRow = "top"
	} else if event.Data.Stage == "production" {
		ufoRow = "bottom"
	}

	if event.Type == "sh.keptn.events.new-artefact" {
		ufoRow := "top"
		ufoColor := "0000ff"
		infoLog.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
		sendUFORequest(ufoAddress, ufoRow, ufoColor, false, true)
	} else if event.Type == "sh.keptn.events.deployment-finished" {
		ufoRow := "top"
		ufoColor := "800080"
		infoLog.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
		sendUFORequest(ufoAddress, ufoRow, ufoColor, false, true)
	} else if event.Type == "sh.keptn.events.tests-finished" {
		ufoRow := "top"
		ufoColor := "00ff00"
		infoLog.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
		sendUFORequest(ufoAddress, ufoRow, ufoColor, true, false)
	} else if event.Type == "sh.keptn.events.evaluation-done" {
		if event.Data.EvaluationPassed {
			ufoRow := "bottom"
			ufoColor := "00ff00"
			infoLog.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
			sendUFORequest(ufoAddress, ufoRow, ufoColor, false, false)
		} else {
			ufoRow := "bottom"
			ufoColor := "ff0000"
			infoLog.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
			sendUFORequest(ufoAddress, ufoRow, ufoColor, false, false)
		}
	}
}

// sendUFORequest : creates and issues necessary GET requests to set UFO LEDs
func sendUFORequest(ufoAddress string, ufoRow string, ufoColor string, morph bool, whirl bool) {
	url := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|15|" + ufoColor
	urlmorph := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|15|" + ufoColor + "&top_morph=30|10"
	urlwhirl := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|1|" + ufoColor + "&top_whirl=240"
	var preparedurl string
	if morph {
		if whirl {
			errorLog.Println("UFO does not support both morphing and whirling at the same time")
			return
		}
		preparedurl = urlmorph
	} else if whirl {
		preparedurl = urlwhirl
	} else {
		preparedurl = url
	}
	req, err := http.NewRequest("GET", preparedurl, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorLog.Println("Error while sending request to UFO: " + err.Error())
		return
	}
	defer resp.Body.Close()
	infoLog.Println("Response Status:" + resp.Status)
}

// ufoInit : upon service start initializes UFO
func ufoInit(ufoAddress string) {
	initURL := "http://" + ufoAddress + "/api?top_init&bottom_init"
	infoLog.Println("Trying to initialize UFO at " + ufoAddress)
	resp, err := http.Get(initURL)
	if err != nil {
		errorLog.Println("Error while sending request to UFO: " + err.Error())
		return
	}
	defer resp.Body.Close()
	infoLog.Println("Response Status:" + resp.Status)
}

func main() {
	Logging(os.Stdout, os.Stderr)
	ufoAddress = os.Getenv("UFO_ADDRESS")
	if ufoAddress == "" {
		errorLog.Println("No UFO address defined")
		return
	}
	ufoInit(ufoAddress)

	http.HandleFunc("/", keptnHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Print("UFO service started.")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
