package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type CreateResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseID      string `json:"responseId"`
	ResponseMessage string `json:"responseMessage"`
	ResponseTime    string `json:"responseTime"`
}

func CreateUser(user UserRequest) (*UserResponse, error) {
	url := "https://b5fwdpwxbl.execute-api.ap-southeast-1.amazonaws.com/v1/api2_1"
	jsonValue, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshal json body: ", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Error create new request: ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		return nil, err
	}
	body := new(UserResponse)
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		log.Println("Error unmarshalling response: ", err)
		return nil, err
	}
	log.Println("Response body when create: ", body)
	return body, nil
}
