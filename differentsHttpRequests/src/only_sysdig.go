package main

import (
	"net/http"

	"io/ioutil"
	"log"
)

func getPrometheusAlerts(url string, token string) {
	client := &http.Client{}

	//Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	var bearer = "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("IBMInstanceID", "4ee1c120-d804-4b54-a0e6-c2ed2364dd63")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Close response body
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(data))
}

func main() {
	getPrometheusAlerts("https://private.eu-de.monitoring.cloud.ibm.com/prometheus/api/v1/alerts", "mytokenlargocifradosuperseguro")
}
