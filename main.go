package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

//Basic HTML Template
var basichtml = `
<html>
<head><title>HTest HTTP Test Tool</title><style>body { font-family: sans-serif; color: #333; padding: 80px;} div {padding: 20px;} h1 { color: #c9510c }</style></head>
<body>
	<div style="background: #e3ecf5; font-size: 20px; font-weight: bold; color: #525252; ">HTest HTTP Test Tool</div>
    <div>%s</div>
	<div style="background: #f3f3f3; text-align: center; border-top: 1px solid #333;">Total Hits: %d | Response Delay: %d | Failure Rate: %d | Jitter: %d</div>
	<div style="color: #525252; text-align: center;">Generated at: %s</div>
</body>
</html>
`

type serverVariables struct {
	Jitter        int // Percent
	FailureRate   int // Percent
	ResponseDelay int // ms
}

var svar serverVariables
var totalHits int

func simpletemplate(content string) string {
	return fmt.Sprintf(basichtml, content, totalHits, svar.ResponseDelay, svar.FailureRate, svar.Jitter, time.Now())
}

func basicrequest(w http.ResponseWriter, r *http.Request) {
	// Really simple hit counter
	totalHits++

	// Display every 100,000 hits
	if math.Mod(float64(totalHits), 100000) == 0 {
		log.Infof("%d requests served", totalHits)
	}

	// Check to see if we need to delay requests
	if svar.ResponseDelay > 0 {
		delay := int64(svar.ResponseDelay)

		// Add jitter (random variance) if requested
		if svar.Jitter > 0 {
			jitter := int64(svar.Jitter)
			delay += delay * (rand.Int63n(jitter*2) - jitter) / 100
		}
		log.Infof("Delay: %v", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	fmt.Fprintf(w, simpletemplate("Default homepage"))
}

func updatevar(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	// Set the variable
	for key := range r.Form {
		switch key {
		case "jitter":
			svar.Jitter, _ = strconv.Atoi(r.FormValue(key))
		case "failurerate":
			svar.FailureRate, _ = strconv.Atoi(r.FormValue(key))
		case "responsedelay":
			svar.ResponseDelay, _ = strconv.Atoi(r.FormValue(key))
		}
	}

	// Possibly redirect here?
	fmt.Fprintf(w, simpletemplate("Value updated"))
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("[HTest] HTTP Test Tool - 1.1 JAN17")

	// Default settings and command line options
	port := flag.Int("port", 8000, "The HTTP port to listen on, default is 8000")
	ipaddr := flag.String("ipaddr", "127.0.0.1", "The IP Address to listen to, default is 127.0.0.1")
	flag.IntVar(&svar.FailureRate, "failurerate", 0, "The amount of requests to return a HTTP 503 instead of 200")
	flag.IntVar(&svar.ResponseDelay, "responsedelay", 0, "The amount to delay requests by if simulating your application response times")
	flag.IntVar(&svar.Jitter, "jitter", 0, "The amount to vary the failurerate and responsedelay by")
	flag.Parse()
	log.Infof("Running on IP %s and port %d", *ipaddr, *port)

	// Start the webserver, with two URL's to handle
	http.HandleFunc("/", basicrequest)
	http.HandleFunc("/svar/", updatevar)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *ipaddr, *port), nil)
	if err != nil {
		log.Panicf("Couldn't listen on IP %s using port %d, check that the details are valid and the port is free", *ipaddr, *port)
	}
}
