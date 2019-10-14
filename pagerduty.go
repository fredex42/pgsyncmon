package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var baseUri = "https://api.pagerduty.com/incidents"

type ServiceReference struct {
	Id         string `json:"id"`
	RecordType string `json:"type"` //must be "service_reference"
}

type PagerDutyBody struct {
	RecordType string `json:"type"` //must be "incident_body"
	Details    string `json:"details"`
}

type PagerDutyIncident struct {
	RecordType  string           `json:"type"` //must be "incident"
	Title       string           `json:"title"`
	Service     ServiceReference `json:"service"`
	Urgency     string           `json:"urgency"`
	IncidentKey string           `json:"incident_key"` //key used to disambiguate incidents.
	Body        PagerDutyBody    `json:"body"`
}

func NewPagerDutyIncident(title string, serviceId string, urgency string, incidentKey string, details string) PagerDutyIncident {
	return PagerDutyIncident{
		RecordType:  "incident",
		Title:       title,
		Service:     ServiceReference{RecordType: "service_reference", Id: serviceId},
		Urgency:     urgency,
		IncidentKey: incidentKey,
		Body:        PagerDutyBody{RecordType: "incident_body", Details: details},
	}
}

func SendIncident(config *Config, description string, incidentKey string, attempt int) error {
	requestBody := NewPagerDutyIncident(config.AlertTitle, config.ServiceId, config.Urgency, incidentKey, description)

	requestContent, marshalErr := json.Marshal(requestBody)
	if marshalErr != nil {
		log.Printf("Could not build requet body: %s", marshalErr)
		return marshalErr
	}

	request, _ := http.NewRequest("POST", baseUri, bytes.NewReader(requestContent))
	request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	request.Header.Set("Authorization", fmt.Sprintf("Token token=%s", config.PagerdutyKey))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not send alert to pagerduty: %s", err)
		return err
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	switch response.StatusCode {
	case 400:
		log.Printf("Invalid arguments. Server said: %s", string(responseBody))
	case 401:
		log.Printf("Authentication was invalid, check your api key is correct. Server said: %s", string(responseBody))
	case 403:
		log.Printf("You API key does not allow you to create incidents, please contact your administrator")
	case 429:
		retryTime := attempt * 2

		log.Printf("Rate limiting prevents this request, retrying in %d seconds...", retryTime)
		time.Sleep(time.Duration(retryTime) * time.Second)
		return SendIncident(config, description, incidentKey, attempt+1)
	case 504:
		fallthrough
	case 503:
		retryTime := attempt * 2

		log.Printf("Pagerduty's server is returning not available, retrying in %d seconds...", retryTime)
		time.Sleep(time.Duration(retryTime) * time.Second)
		return SendIncident(config, description, incidentKey, attempt+1)
	case 200:
		log.Printf("Incident sent to PD")
	case 201:
		log.Printf("Incident sent to PD")
	default:
		log.Printf("Got unexpected status code %d from server, server said: %s", response.StatusCode, string(responseBody))
	}

	return nil
}
