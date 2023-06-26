// Import necessary packages
package main

import (
	"crypto/hmac"   // Package to compute message authentication codes using HMAC (Keyed-Hashing for Message Authentication)
	"crypto/sha256" // Package to implement the SHA-256 hashing algorithm
	"encoding/hex"  // Package to encode/decode hexadecimal strings
	"log"           // Package to log messages
	"time"          // Package to work with time objects

	"github.com/aws/aws-lambda-go/events"            // AWS Lambda package for Go events
	"github.com/aws/aws-lambda-go/lambda"            // AWS Lambda package for Go
	"github.com/awslabs/aws-lambda-go-api-proxy/gin" // AWS Lambda custom runtime with Gin (web framework) support
	"github.com/gin-gonic/gin"                       // Gin is a web framework for Go
)

// Declare global variable for ginLambda (AWS Lambda with Gin)
var ginLambda *ginadapter.GinLambda

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
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(plaintext))
	return hex.EncodeToString(h.Sum(nil))
}

// Initialize Gin application and setup routes
func init() {
	log.Print("Gin cold start")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Route handler for POST requests to /api2
	r.POST("/api2", func(c *gin.Context) {
		var reqBody RequestBody
		c.BindJSON(&reqBody) // Bind request JSON to reqBody struct
		log.Printf("%+v", reqBody)

		// Compute HMAC-SHA256 signature of plainText using secretKey and get current time
		signature := hmacSha256(reqBody.Data.PlainText, reqBody.Data.SecretKey)
		currentTime := time.Now().Format(time.RFC3339)

		// Send response JSON with requestId, requestTime, and computed signature
		c.JSON(200, gin.H{
			"requestId":   reqBody.RequestId,
			"requestTime": currentTime,
			"data": gin.H{
				"signature": signature,
			},
		})
	})

	// Initialize ginLambda with the Gin routes
	ginLambda = ginadapter.New(r)
}

// Handler function for AWS Lambda
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.Proxy(req)
}

// Main function for AWS Lambda
func main() {
	lambda.Start(Handler)
}
