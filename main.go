package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	keptnevents "github.com/akirasoft/go-utils/pkg/events"
	keptnutils "github.com/akirasoft/go-utils/pkg/utils"
	"github.com/kelseyhightower/envconfig"
)

var ufoAddress string

var ufoRow string

// setUfoRow sets ufoRow to bottom when stage starts with prod and top when stage starts with dev or stag, otherwise defaults to top.
// Supports longer form stage definitions that might be in the shipyard file
func setUfoRow(stage string) string {
	stageUpper := strings.ToUpper(stage)
	if strings.HasPrefix(stageUpper, "DEV") {
		ufoRow = "top"
	} else if strings.HasPrefix(stageUpper, "STAG") {
		ufoRow = "top"
	} else if strings.HasPrefix(stageUpper, "PROD") {
		ufoRow = "bottom"
	} else {
		ufoRow = "top"
	}
	//log.Println("UfoRow will be:", ufoRow)
	return ufoRow
}

// ufoReceiver receives keptn events via http and sets UFO LEDs based on payload
func ufoReceiver(data interface{}, shkeptncontext string, eventID string) error {
	logger := keptnutils.NewLogger(shkeptncontext, eventID, "ufo-service")
	switch data.(type) {
	case *keptnevents.EvaluationDoneEvent:
		var event = data.(*keptnevents.EvaluationDoneEvent)
		ufoRow := setUfoRow(event.Stage)
		if event.Evaluationpassed {
			ufoColor := "00ff00"
			logger.Info(fmt.Sprintln("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor))
			go sendUFORequest(ufoAddress, ufoRow, ufoColor, false, false, logger)
		} else {
			ufoColor := "ff0000"
			logger.Info(fmt.Sprintln("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor))
			go sendUFORequest(ufoAddress, ufoRow, ufoColor, false, false, logger)
		}
	case *keptnevents.ConfigurationChanged:
		var event = data.(*keptnevents.ConfigurationChanged)
		ufoRow := setUfoRow(event.Stage)
		ufoColor := "0000ff"
		logger.Info(fmt.Sprintln("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor))
		go sendUFORequest(ufoAddress, ufoRow, ufoColor, false, true, logger)
	case *keptnevents.DeploymentFinishedEvent:
		var event = data.(*keptnevents.DeploymentFinishedEvent)
		ufoRow := setUfoRow(event.Stage)
		ufoColor := "800080"
		logger.Info(fmt.Sprintln("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor))
		go sendUFORequest(ufoAddress, ufoRow, ufoColor, false, true, logger)
	case *keptnevents.TestsFinishedEvent:
		var event = data.(*keptnevents.TestsFinishedEvent)
		ufoRow := setUfoRow(event.Stage)
		ufoColor := "00ff00"
		logger.Info(fmt.Sprintln("Trying to talk to UFO at " + ufoAddress + " setting " + ufoRow + " to " + ufoColor))
		go sendUFORequest(ufoAddress, ufoRow, ufoColor, true, false, logger)
	default:
		logger.Info("Other event")
	}

	return nil
}

func sendUFORequest(ufoAddress string, ufoRow string, ufoColor string, morph bool, whirl bool, logger *keptnutils.Logger) {
	url := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|15|" + ufoColor
	urlmorph := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|15|" + ufoColor + "&" + ufoRow + "_morph=30|10"
	urlwhirl := "http://" + ufoAddress + "/api?" + ufoRow + "_init&" + ufoRow + "=0|1|" + ufoColor + "&" + ufoRow + "_whirl=240"
	var preparedurl string
	if morph {
		if whirl {
			logger.Error("UFO does not support both morphing and whirling at the same time")
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
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error while sending request to UFO: " + err.Error())
		return
	}
	defer resp.Body.Close()
	logger.Info(fmt.Sprintln("Response Status:" + resp.Status))
}

func main() {
	ufoAddress = os.Getenv("UFO_ADDRESS")
	if ufoAddress == "" {
		log.Println("No UFO address defined")
		os.Exit(1)
	}

	var rcv keptnutils.RcvConfig
	if err := envconfig.Process("", &rcv); err != nil {
		log.Printf("[ERROR] Failed to process listener var: %s", err)
		os.Exit(1)
	}

	keptnutils.KeptnReceiver(rcv, ufoReceiver)

}
