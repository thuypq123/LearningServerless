package common

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type RequestBody struct {
	RequestID string `json:"requestId"`
	Data      Data   `json:"data"`
}
type Data struct {
	Value int `json:"value"`
}
type ResponseBody struct {
	ResponseCode    string `json:"responseCode"`
	ResponseID      string `json:"responseId"`
	ResponseMessage string `json:"responseMessage"`
	ResponseTime    string `json:"responseTime"`
}

func VerifyPhone(phone int) (bool, error) {
	url := "https://1g1zcrwqhj.execute-api.ap-southeast-1.amazonaws.com/dev/testapi"
	requestBody := &RequestBody{
		RequestID: uuid.New().String(),
		Data: Data{
			Value: phone,
		},
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error marshal json body: ", err)
		return false, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Error create new request: ", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request: ", err)
		return false, err
	}
	defer resp.Body.Close()
	body := new(ResponseBody)
	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		log.Println("Error decoding response body [verifyPhone]: ", err)
		return false, err
	}
	if body.ResponseCode == "00" {
		return true, nil
	}
	return false, nil
}
