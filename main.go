package main

import (
	"log"
	"net/http"
	"os"

	keptnevents "github.com/akirasoft/keptn-events"
	"github.com/kelseyhightower/envconfig"
)

var ufoAddress string
var ufoRow string

// ufoReceiver receives keptn events via http and sets UFO LEDs based on payload
func ufoReceiver(data interface{}) error {
	switch data.(type) {
	case *keptnevents.EvaluationDoneEvent:
		var event = data.(*keptnevents.EvaluationDoneEvent)
		if event.Stage == "dev" {
			ufoRow = "top"
		} else if event.Stage == "staging" {
			ufoRow = "top"
		} else if event.Stage == "production" {
			ufoRow = "bottom"
		}
		if event.Evaluationpassed {
			ufoColor := "00ff00"
			log.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
			sendUFORequest(ufoAddress, ufoRow, ufoColor, false, false)
		} else {
			ufoColor := "ff0000"
			log.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
			sendUFORequest(ufoAddress, ufoRow, ufoColor, false, false)
		}
	case *keptnevents.NewArtifactEvent:
		var event = data.(*keptnevents.NewArtifactEvent)
		if event.Stage == "dev" {
			ufoRow = "top"
		} else if event.Stage == "staging" {
			ufoRow = "top"
		} else if event.Stage == "production" {
			ufoRow = "bottom"
		}
		ufoColor := "0000ff"
		log.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
		sendUFORequest(ufoAddress, ufoRow, ufoColor, false, true)
	case *keptnevents.DeploymentFinishedEvent:
		var event = data.(*keptnevents.DeploymentFinishedEvent)
		if event.Stage == "dev" {
			ufoRow = "top"
		} else if event.Stage == "staging" {
			ufoRow = "top"
		} else if event.Stage == "production" {
			ufoRow = "bottom"
		}
		ufoColor := "800080"
		log.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
		sendUFORequest(ufoAddress, ufoRow, ufoColor, false, true)
	case *keptnevents.TestsFinishedEvent:
		var event = data.(*keptnevents.TestsFinishedEvent)
		if event.Stage == "dev" {
			ufoRow = "top"
		} else if event.Stage == "staging" {
			ufoRow = "top"
		} else if event.Stage == "production" {
			ufoRow = "bottom"
		}
		ufoColor := "00ff00"
		log.Println("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor)
		sendUFORequest(ufoAddress, ufoRow, ufoColor, true, false)
	default:
		log.Println("Other event")
	}

	return nil
}

func sendUFORequest(ufoAddress string, ufoRow string, ufoColor string, morph bool, whirl bool) {
	url := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|15|" + ufoColor
	urlmorph := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|15|" + ufoColor + "&" + ufoRow + "_morph=30|10"
	urlwhirl := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|1|" + ufoColor + "&" + ufoRow + "_whirl=240"
	var preparedurl string
	if morph {
		if whirl {
			log.Println("UFO does not support both morphing and whirling at the same time")
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
		log.Println("Error while sending request to UFO: " + err.Error())
		return
	}
	defer resp.Body.Close()
	log.Println("Response Status:" + resp.Status)
}

func main() {
	ufoAddress = os.Getenv("UFO_ADDRESS")
	if ufoAddress == "" {
		log.Println("No UFO address defined")
		os.Exit(1)
	}

	var rcv keptnevents.RcvConfig
	if err := envconfig.Process("", &rcv); err != nil {
		log.Printf("[ERROR] Failed to process listener var: %s", err)
		os.Exit(1)
	}

	keptnevents.KeptnReceiver(rcv, ufoReceiver)

}
