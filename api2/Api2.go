// Import necessary packages
package main

import (
	"crypto/hmac"   // Package to compute message authentication codes using HMAC (Keyed-Hashing for Message Authentication)
	"crypto/sha256" // Package to implement the SHA-256 hashing algorithm
	"encoding/hex"  // Package to encode/decode hexadecimal strings
	"encoding/json" // Package to deal with JSON data format
	"log"           // Package to log messages
	"net/http"      // Package to build HTTP server
	"time"          // Package to work with time objects

	"github.com/aws/aws-lambda-go/events" // AWS Lambda package for Go events
	"github.com/aws/aws-lambda-go/lambda" // AWS Lambda package for Go
)

// Define RequestBody structure to parse request JSON
type RequestBody struct {
	RequestId   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Data        struct {
		PlainText string `json:"plainText"`
		SecretKey string `json:"secretKey"`
	} `json:"data"`
}

// hmacSha256 computes the HMAC-SHA256 signature of plaintext using secretKey
func hmacSha256(plaintext string, secretKey string) string {
	log.Println("Computing HMAC-SHA256 signature of plaintext using secretKey")
	log.Println("plaintext: ", plaintext)
	log.Println("secretKey: ", secretKey)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(plaintext))
	return hex.EncodeToString(h.Sum(nil))
}

// Handler function for AWS Lambda
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse JSON from request body
	var reqBody RequestBody
	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		// Return status code 400 (Bad Request) if JSON parsing fails
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request",
		}, nil
	}
	log.Printf("%+v", reqBody)

	// Compute HMAC-SHA256 signature of plainText using secretKey and get current time
	signature := hmacSha256(reqBody.Data.PlainText, reqBody.Data.SecretKey)
	currentTime := time.Now().Format(time.RFC3339)

	// Create response body
	respBody := map[string]interface{}{
		"requestId":   reqBody.RequestId,
		"requestTime": currentTime,
		"data": map[string]string{
			"signature": signature,
		},
	}
	respJson, _ := json.Marshal(respBody)

	// Return status code 200 (OK) and response JSON
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respJson),
	}, nil
}

// Main function for AWS Lambda
func main() {
	lambda.Start(Handler)
}
