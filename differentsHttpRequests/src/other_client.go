package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"fmt"
	"io/ioutil"
	"log"
)

type VaultTokenResponse struct {
	RequestID     string   `json:"request_id"`
	LeaseID       string   `json:"lease_id"`
	Renewable     bool     `json:"renewable"`
	LeaseDuration int      `json:"lease_duration"`
	Data          string   `json:"data"`
	WrapInfo      string   `json:"wrap_info"`
	Warnings      []string `json:"warnings"`
	Auth          struct {
		ClientToken   string   `json:"client_token"`
		Accessor      string   `json:"accessor"`
		Policies      []string `json:"policies"`
		TokenPolicies []string `json:"token_policies"`
		Metadata      struct {
			AuthorityKeyID string `json:"authority_key_id"`
			CertName       string `json:"cert_name"`
			CommonName     string `json:"common_name"`
			SerialNumber   string `json:"serial_number"`
			SubjectKeyID   string `json:"subject_key_id"`
		} `json:"metadata"`
		LeaseDuration  int    `json:"lease_duration"`
		Renewable      bool   `json:"renewable"`
		EntityID       string `json:"entity_id"`
		TokenType      string `json:"token_type"`
		Orphan         bool   `json:"orphan"`
		MfaRequirement string `json:"mfa_requirement"`
		NumUses        int    `json:"num_uses"`
	} `json:"auth"`
}

type VaultApiResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		APIKey string `json:"api_key"`
	} `json:"data"`
	WrapInfo string `json:"wrap_info"`
	Warnings string `json:"warnings"`
	Auth     string `json:"auth"`
}

type IbmTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Expiration   int    `json:"expiration"`
	Scope        string `json:"scope"`
}

func getVaultToken(url string, body string) (string, error) {

	// Load client cert
	certFile := "/opt/certificate/tls.crt"
	keyFile := "/opt/certificate/tls.key"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caFile := "/opt/ca/ca-certificates.crt"
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Prepare a JSON body
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(bodyBytes)

	// Make HTTP POST request (with Do)
	request, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return "error creating http post request", err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		return "error making http post request", err
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

	// Body response to json
	var result VaultTokenResponse
	if err := json.Unmarshal(responseBody, &result); err != nil { // Parse []byte to go struct
		fmt.Println("Can not unmarshal JSON")
	}
	log.Println("Token: ", string(result.Auth.ClientToken))
	return result.Auth.ClientToken, nil
}

func getIbmApiKey(url string, token string) (string, error) {
	// Load CA cert
	caFile := "/opt/ca/ca-certificates.crt"
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

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

	// Body response to json
	var result VaultApiResponse
	if err := json.Unmarshal(data, &result); err != nil { // Parse []byte to go struct
		fmt.Println("Can not unmarshal JSON")
	}
	log.Println("Token: ", string(result.Data.APIKey))

	return result.Data.APIKey, nil
}

func getIbmToken(url string, body string) (string, error) {

	client := &http.Client{}

	// Prepare body
	var data = strings.NewReader(body)

	// Make HTTP POST request (with Do)
	request, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		return "", err
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

	// Body response to json
	var result IbmTokenResponse
	if err := json.Unmarshal(responseBody, &result); err != nil { // Parse []byte to go struct
		fmt.Println("Can not unmarshal JSON")
	}
	log.Println("Token: ", string(result.AccessToken))

	return result.AccessToken, nil
}

func main() {
	token, err := getVaultToken("https://kvault.cloud.caixabank.com/v1/auth/cert-vault-auth/login", "{\"name\": \"pcld-bchicp-p-pro\"}")
	if err != nil {
		log.Fatal(err)
	}
	apikey, err := getIbmApiKey("https://kvault.cloud.caixabank.com/v1/ibmcloud/cxb-ope-pro/creds/pbchicp-icpco-rg01-pro", token)
	if err != nil {
		log.Fatal(err)
	}
	ibmOAuth := "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=" + apikey
	ibmToken, err := getIbmToken("https://private.iam.cloud.ibm.com/identity/token", ibmOAuth)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Token (en el main): ", string(ibmToken))
}
