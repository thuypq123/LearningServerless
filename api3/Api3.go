package main

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

// Declare a global variable for GinLambda adapter
var ginLambda *ginadapter.GinLambda

// Define the request body structure
type RequestBody struct {
	RequestId   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Data        struct {
		NeedEncode string `json:"needEncode"`
		NeedDecode string `json:"needDecode"`
	} `json:"data"`
}

// Initialize the Gin router and Lambda adapter
func init() {
	log.Print("Gin cold start")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Setup POST route for encoding and decoding base64
	r.POST("/api3", func(c *gin.Context) {
		// Bind JSON request to RequestBody struct
		var reqBody RequestBody
		c.BindJSON(&reqBody)

		// Log the request body
		log.Printf("%+v", reqBody)

		// Encode the data
		encodedData := base64.StdEncoding.EncodeToString([]byte(reqBody.Data.NeedEncode))
		log.Printf("Encoded data: %s", encodedData)

		// Decode the data
		decodedData, err := base64.StdEncoding.DecodeString(reqBody.Data.NeedDecode)
		log.Printf("Decoded data: %s", decodedData)

		if err != nil {
			// Log decoding error
			log.Printf("Error decoding base64: %s", err)
		} else {
			// Log decoded data
			log.Printf("Decoded data: %s", decodedData)
		}

		// Get the current time
		currentTime := time.Now().Format(time.RFC3339)

		// Send JSON response
		c.JSON(200, gin.H{
			"requestId":   reqBody.RequestId,
			"requestTime": currentTime,
			"data": gin.H{
				"outEncode": encodedData,
				"outDecode": string(decodedData),
			},
		})
	})

	// Create a new GinLambda adapter with the Gin router
	ginLambda = ginadapter.New(r)
}

// Handler function for AWS Lambda
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.Proxy(req)
}

// Main function to start the Lambda function
func main() {
	lambda.Start(Handler)
}
