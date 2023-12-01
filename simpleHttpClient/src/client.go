package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"io/ioutil"
	"log"
)

//var (
//	certFile = flag.String("cert", "routeToSomeCertFile", "A PEM encoded certificate file.")
//	keyFile  = flag.String("key", "routeToSomeKeyFile", "A PEM encoded private key file.")
//	caFile   = flag.String("CA", "routeToSomeCertCAFile", "A PEM eoncoded CA's certificate file.")
//	//service  = flag.String("url", "https://localhost:8080", "The service URL which the request will be made")
//)
//
//func makeGetTLSRequest() {
//	flag.Parse()
//
//	// Load client cert
//	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Load CA cert
//	caCert, err := ioutil.ReadFile(*caFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//	caCertPool := x509.NewCertPool()
//	caCertPool.AppendCertsFromPEM(caCert)
//
//	// Setup HTTPS client
//	tlsConfig := &tls.Config{
//		Certificates: []tls.Certificate{cert},
//		RootCAs:      caCertPool,
//	}
//	tlsConfig.BuildNameToCertificate()
//	transport := &http.Transport{TLSClientConfig: tlsConfig}
//	client := &http.Client{Transport: transport}
//
//	// Do GET something
//	resp, err := client.Get("https://goldportugal.local:8443")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	// Dump response
//	data, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println(string(data))
//}
//
func makeGetWithoutTLS(url string, token string) {
	client := &http.Client{}

	//Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set header
	req.Header.Set("X-Vault-Token", token)

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

func makePostRequest(url string, body string) {

	// Prepare a JSON body
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(bodyBytes)

	// Make HTTP POST request
	//resp, err := http.Post(url, "application/json", reader)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// Make HTTP POST request (with Do)
	request, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		return
	}

	// Close response body
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		log.Println("Error response. Status Code: ", resp.StatusCode)
	}

	log.Println("Response:", string(responseBody))
}

func main() {
	makePostRequest("http://0.0.0.0:9000", "{\"field1\": \"value1\"}")
	makeGetWithoutTLS("http://0.0.0.0:9000", "mytokenlargocifradosuperseguro")
}
