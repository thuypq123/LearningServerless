package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestBody struct {
	RequestID string `json:"requestId"`
	Data      struct {
		Value int `json:"value"`
	} `json:"data"`
}

type ResponseBody struct {
	ResponseCode    string `json:"responseCode"`
	ResponseID      string `json:"responseId"`
	ResponseMessage string `json:"responseMessage"`
	ResponseTime    string `json:"responseTime"`
}

func callAPI(value int) (*ResponseBody, error) {
	requestBody := &RequestBody{
		RequestID: "13456",
		Data: struct {
			Value int `json:"value"`
		}{
			Value: value,
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://1g1zcrwqhj.execute-api.ap-southeast-1.amazonaws.com/dev/testapi", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseBody := new(ResponseBody)
	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return responseBody, nil
}

func main() {
	resp, err := callAPI(1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(resp)
}
