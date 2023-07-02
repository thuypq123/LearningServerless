package common

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Struct to model the data part of the request body
type Data struct {
	Value int `json:"value"`
}

// Struct to model the entire request body.
type RequestBody struct {
	RequestID string `json:"requestId"` // Unique identifier for the request
	Data      Data   `json:"data"`      // Payload data for the request
}

// Struct to model the response body that comes back from the API.
type ResponseBody struct {
	ResponseCode    string `json:"responseCode"`    // The status code of the response
	ResponseID      string `json:"responseId"`      // Unique identifier for the response
	ResponseMessage string `json:"responseMessage"` // Any message sent in the response
	ResponseTime    string `json:"responseTime"`    // Timestamp of the response
}

// Function to verify the phone number by sending a POST request to the provided API.
func VerifyPhone(phone int) (bool, error) {
	url := "https://1g1zcrwqhj.execute-api.ap-southeast-1.amazonaws.com/dev/testapi"

	// Create the request body using the provided phone number and a new UUID.
	requestBody := &RequestBody{
		RequestID: uuid.New().String(),
		Data: Data{
			Value: phone,
		},
	}

	// Convert the request body into JSON format
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error marshal json body: ", err)
		return false, err
	}

	// Create a new HTTP request with the JSON body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Error create new request: ", err)
		return false, err
	}

	// Set the headers for the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request: ", err)
		return false, err
	}
	defer resp.Body.Close()

	// Create a new ResponseBody struct to hold the response
	body := new(ResponseBody)

	// Decode the JSON response into the ResponseBody struct
	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		log.Println("Error decoding response body [verifyPhone]: ", err)
		return false, err
	}

	// Return true if the ResponseCode is "00", indicating a successful verification
	if body.ResponseCode == "00" {
		return true, nil
	}

	// Otherwise return false
	return false, nil
}
