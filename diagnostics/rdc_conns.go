package diagnostics

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"path"
    "strconv"


	log "github.com/sirupsen/logrus"
)

// RDCServices makes connections to the main RDC endpoints to prove
// that the endpoints are reachable from the machine
func RDCServices(rdcEndpoints []string) {
	for _, endpoint := range rdcEndpoints {
		log.Debug("Sending POST req to ", endpoint)
		var jsonBody = []byte(`{"test":"this will result in an HTTP 500 resp or 401 resp."}`)
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[ ] %s not reachable\n", endpoint)
			log.WithFields(log.Fields{
				"error": err,
			}).Infof("[ ] %s not reachable\n", endpoint)
		}

		if err == nil {
			respOutput(resp, endpoint)
		}
	}
}

func _SendRequest(endpoint string, timeout int) (*http.Response, error) {
	client := &http.Client{Timeout: time.Duration(timeout + 2) * time.Second}
	u, _ := url.Parse(endpoint)
	u.Path = path.Join(u.Path, strconv.Itoa(timeout))
	endpoint = u.String()
	log.Info("Sending GET request to ", endpoint)
	resp, err := client.Get(endpoint)
	return resp, err
}

func LongIdleConnections(endpoint string) {
	got_reply := false
	// seconds := 15 * 60 // 15 minutes
	lowerBound := 0
	upperBound := 130
	seconds := upperBound
	log.Debug("Initial timeout is ", seconds, " seconds.")
	_, err := url.Parse(endpoint) // Doesn't work?
	if err != nil {
		log.Error("Malformed endpoint: ", endpoint)
		return
	}

	for {
		start := time.Now()
		resp, err := _SendRequest(endpoint, seconds)

		elapsed := time.Since(start)
		if resp != nil {
			log.Info("Got a reply after ", elapsed, " seconds.")

			if seconds == upperBound || (upperBound - lowerBound) < 10 {
				got_reply = true
			} else {
				lowerBound = seconds
				seconds = seconds + ((upperBound - seconds) / 2)
				log.Debug("Changing timeout to ", seconds, " seconds.")
			}
		} else if err != nil {
			log.Error("Didn't receive any reply after ", elapsed, " seconds: request cancelled.")
			log.Error(err)
			seconds = lowerBound + (upperBound - lowerBound) / 2.0
			if (upperBound - seconds) < 10 {
				log.Debug("Approximation is less than 10 seconds: Abort.")
				break
			}
			log.Debug("Changing timeout to ", seconds, " seconds.")
		}

		if got_reply {
			break
		}
	}

	log.Info("Approximate time for the longuest idle connection possible is ", lowerBound, " seconds.")
}
