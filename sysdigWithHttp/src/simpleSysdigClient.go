package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type SysdigAlertsResponse struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []struct {
			Labels struct {
				Severity           string `json:"severity"`
				Alertname          string `json:"alertname"`
				KubeNamespaceName  string `json:"kube_namespace_name,omitempty"`
				KubeDeploymentName string `json:"kube_deployment_name,omitempty"`
				KubeClusterName    string `json:"kube_cluster_name,omitempty"`
				KubeNodeName       string `json:"kube_node_name,omitempty"`
				KubePodName        string `json:"kube_pod_name,omitempty"`
			} `json:"labels,omitempty"`
			Annotations struct {
				Description string `json:"description"`
			} `json:"annotations"`
			State    string  `json:"state"`
			ActiveAt string  `json:"activeAt"`
			Value    float64 `json:"value"`
		} `json:"alerts"`
	} `json:"data"`
}

//func getPrometheusAlerts(url string, token string) []SysdigAlertsResponse {
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

	//log.Println(string(data))
	// Body response to json
	var result SysdigAlertsResponse
	//if err := json.Unmarshal(data, &result); err != nil { // Parse []byte to go struct
	if err := json.Unmarshal([]byte(data), &result); err != nil { // Parse []byte to go struct
		log.Fatal("Can not unmarshal JSON")
	}

	//log.Println(result)
	//log.Println("Inicio")
	//log.Println(result.Data)
	//log.Println("Fin")
	for _, item := range result.Data.Alerts {
		log.Println("State --> ", item.State)
		log.Println("Serveriry --> ", item.Labels.Severity)
		log.Println("Labels, Alertname --> ", item.Labels.Alertname)
		log.Println("Value: ", item.Value)
		log.Println("ActiveAt: ", item.ActiveAt)
		log.Println("Cluster: ", item.Labels.KubeClusterName)
		log.Println("---Fin---")
	}
	//return result.Data.Alerts
}

func main() {
	//getPrometheusAlerts("https://private.eu-de.monitoring.cloud.ibm.com/prometheus/api/v1/alerts", "mytokenlargocifradosuperseguro")
	getPrometheusAlerts("https://private.eu-de.monitoring.cloud.ibm.com/prometheus/api/v1/alerts", "eyJraWQiOiIyMDIzMTIwNzA4MzYiLCJhbGciOiJSUzI1NiJ9.eyJpYW1faWQiOiJpYW0tU2VydmljZUlkLWU2ODhmMWIwLWZjMzUtNGJlZC1iZDM0LTkyYzVjMjc2M2FhNyIsImlkIjoiaWFtLVNlcnZpY2VJZC1lNjg4ZjFiMC1mYzM1LTRiZWQtYmQzNC05MmM1YzI3NjNhYTciLCJyZWFsbWlkIjoiaWFtIiwianRpIjoiZWVmNDI1NjYtZTA5Yi00OTcxLThiODEtMDU4OTI5ZmYzYzg5IiwiaWRlbnRpZmllciI6IlNlcnZpY2VJZC1lNjg4ZjFiMC1mYzM1LTRiZWQtYmQzNC05MmM1YzI3NjNhYTciLCJuYW1lIjoidmF1bHQtZ2VuZXJhdGVkLXBiY2hpY3AtaWNwY28tcmcwMS1wcm8iLCJzdWIiOiJTZXJ2aWNlSWQtZTY4OGYxYjAtZmMzNS00YmVkLWJkMzQtOTJjNWMyNzYzYWE3Iiwic3ViX3R5cGUiOiJTZXJ2aWNlSWQiLCJhdXRobiI6eyJzdWIiOiJTZXJ2aWNlSWQtZTY4OGYxYjAtZmMzNS00YmVkLWJkMzQtOTJjNWMyNzYzYWE3IiwiaWFtX2lkIjoiaWFtLVNlcnZpY2VJZC1lNjg4ZjFiMC1mYzM1LTRiZWQtYmQzNC05MmM1YzI3NjNhYTciLCJzdWJfdHlwZSI6IlNlcnZpY2VJZCIsIm5hbWUiOiJ2YXVsdC1nZW5lcmF0ZWQtcGJjaGljcC1pY3Bjby1yZzAxLXBybyJ9LCJhY2NvdW50Ijp7InZhbGlkIjp0cnVlLCJic3MiOiJmMmYxYWZkMDJkNTc0YTMzYWJiZjg2NDU3YjViNTE0YiIsImZyb3plbiI6dHJ1ZX0sImlhdCI6MTcwMjI5NjQ4NywiZXhwIjoxNzAyMzAwMDg3LCJpc3MiOiJodHRwczovL2lhbS5jbG91ZC5pYm0uY29tL2lkZW50aXR5IiwiZ3JhbnRfdHlwZSI6InVybjppYm06cGFyYW1zOm9hdXRoOmdyYW50LXR5cGU6YXBpa2V5Iiwic2NvcGUiOiJpYm0gb3BlbmlkIiwiY2xpZW50X2lkIjoiZGVmYXVsdCIsImFjciI6MSwiYW1yIjpbInB3ZCJdfQ.cxeYb4GEdPZY__UqVcsPDtvXOntO29YGDMBapk7K_jJJDYR74jZbd2WUtnqfcCr1wixXK00iY6IQm3pa9ydDIBgKKpy4eX24LiSFdiSR8UII4KmrB1hoD8njENQjU1XWPglKe6-1EM5izN-Z7VeI8AJAKWVTbkA9bRxwwwAWZ3wgHDbCrWg-ASnllblvujzjREPCSACPd39pwbN5SsUA2PvJLdkPachfqrsIIBzO0KINTZpkWklCStbUwKXntTzpFMLePHgX4EQxgy4Me8rCV9Z4iqr2q36uAA-L4d5YoiH_sCQlNjojJVyEoc3G_koj6meup4kWj8MFfBppEIb4QA")
	//log.Println(alerts)
}
