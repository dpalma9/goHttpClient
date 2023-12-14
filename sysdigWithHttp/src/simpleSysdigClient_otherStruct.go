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
		Alerts []SysdigAlerts
	} `json:"data"`
}

type SysdigAlerts struct {
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
}

//func getPrometheusAlerts(url string, token string) []SysdigAlertsResponse {
func getPrometheusAlerts(url string, token string) []SysdigAlerts {
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

	//var resultado SysdigAlerts
	resultado := []SysdigAlerts(result.Data.Alerts)
	return resultado
}

func main() {
	alerts := getPrometheusAlerts("https://private.eu-de.monitoring.cloud.ibm.com/prometheus/api/v1/alerts", "mytokenlargocifradosuperseguro")
	log.Println(alerts)
}
