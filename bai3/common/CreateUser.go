package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// CreateResponse struct defines a structure to hold the response from the API after user creation.
type CreateResponse struct {
	ResponseCode    string `json:"responseCode"`    // Holds the status code of the response
	ResponseID      string `json:"responseId"`      // Holds the ID of the response
	ResponseMessage string `json:"responseMessage"` // Holds the response message
	ResponseTime    string `json:"responseTime"`    // Holds the time the response was generated
}

// CreateUser function takes in a user request, sends it to an API to create the user and receives the response.
func CreateUser(user UserRequest) (*UserResponse, error) {
	url := "https://b5fwdpwxbl.execute-api.ap-southeast-1.amazonaws.com/v1/api2_1" // URL of the API endpoint

	// Marshal the user data to JSON format
	jsonValue, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshal json body: ", err)
		return nil, err
	}
	// Creating new HTTP request for API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Error creating new request: ", err)
		return nil, err
	}

	// Setting Content-Type as application/json for request header
	req.Header.Set("Content-Type", "application/json")

	// Using default HTTP Client to send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request: ", err)
		return nil, err
	}
	defer resp.Body.Close() // Closing the response body once done using it

	// Reading all the bytes of the response Body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		return nil, err
	}

	// Unmarshalling the response bytes into UserResponse struct
	body := new(UserResponse)
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		log.Println("Error unmarshalling response: ", err)
		return nil, err
	}

	log.Println("Response body when create: ", body) // Logging the received response body

	// Returning the response body for further usage
	return body, nil
}
