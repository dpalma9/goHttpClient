package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
)

type userAgentRoundTripper struct {
	ibmInstanceId string
	rt            http.RoundTripper
}

// NewUserAgentRoundTripper adds the user agent every request header.
func NewUserAgentRoundTripper(ibmInstanceId string, rt http.RoundTripper) http.RoundTripper {
	return &userAgentRoundTripper{ibmInstanceId, rt}
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	//req = cloneRequest(req)
	req.Header.Set("IBMInstanceID", rt.ibmInstanceId)
	return rt.rt.RoundTrip(req)
}

func GetSysdigHttpTransport() *http.Transport {

	// la parte del time
	timeOut := time.Second
	tOut, err := strconv.ParseFloat("30", 64)
	if err != nil {
		log.Println("Error converting timeOut value", err)
	}
	timeOut = time.Duration(tOut) * time.Second

	tlsHandshakeTimeOut := time.Second
	tTlsOut, err := strconv.ParseFloat("10", 64)
	if err != nil {
		log.Println("Error converting timeOut value", err)
	}
	tlsHandshakeTimeOut = time.Duration(tTlsOut) * time.Second

	transport := &http.Transport{
		//transport := &CustomTransport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeOut,
			KeepAlive: timeOut,
		}).DialContext,
		TLSHandshakeTimeout: tlsHandshakeTimeOut,
	}

	return transport
}

// To get Prometheus (Sysdig) client
func GetSysdigClient() (api.Client, error) {

	log.Println("Preparing Prometheus (sysdig) Client")

	tr := GetSysdigHttpTransport()

	log.Println("Getting Prometheus (sysdig) Client: Getting token")
	api_token := "mitoken"

	urlEndpoint := "https://private.eu-de.monitoring.cloud.ibm.com"
	log.Println("URL to Prometheus queries is: ", urlEndpoint)

	rt := config.NewAuthorizationCredentialsRoundTripper("Bearer", config.Secret(api_token), tr)
	ibmId := "4ee1c120-d804-4b54-a0e6-c2ed2364dd63"
	rt = NewUserAgentRoundTripper(ibmId, rt)

	client, err := api.NewClient(api.Config{
		Address:      urlEndpoint,
		RoundTripper: rt,
	})
	if err != nil {
		return client, err
	}

	return client, nil
}

func main() {
	log.Println("Empieza el main")
	client, err := GetSysdigClient()
	if err != nil {
		log.Fatal("Error al iniciar el cliente: ", err)
	}
	apiPrometheus := v1.NewAPI(client)
	log.Println("Preparing context with client config")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	log.Println("Getting Prometheus Alerts")
	resultPrometheus, err := apiPrometheus.Alerts(ctx)
	if err != nil {
		log.Fatal("Error Getting alerts: ", err)
	}
	log.Println("Estos son las alertas: ")
	log.Println(resultPrometheus)
}
