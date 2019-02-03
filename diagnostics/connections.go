package diagnostics

import (
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

//PublicSites attempts to send HTTP requests to sites that SHOULD be reachable.
func PublicSites(sitelist []string) {
	for _, site := range sitelist {
		log.Debug("Sending GET req to ", site)
		resp, err := http.Get(site)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"resp":  resp,
			}).Errorf("[ ] %s not reachable %s\n", site, resp.Status)
		}

		if resp.StatusCode == 200 {
			fmt.Printf("[\u2713] %s is reachable %s\n", site, resp.Status)
		}
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[\u2713] %s reachable.\n", site)
	}
}

//TCPConns attempts to open various TCP connections to the provided sites
func TCPConns(sitelist []string) {
	for _, site := range sitelist {
		timeout := time.Duration(5 * time.Second)
		conn, err := net.DialTimeout("tcp4", site, timeout)
		if err != nil {
			log.Errorf("%s unreachable, %v: ", site, err)
		}
		fmt.Println("[\u2713] TCP (IPv4) connection to", site)
		log.WithFields(log.Fields{
			"local":  conn.LocalAddr(),
			"remote": conn.RemoteAddr(),
		}).Infof("[\u2713] %s reachable via TCP (IPv4).\n", site)
	}
}