// Import necessary packages
package main

import (
	"log"  // Package to log messages
	"time" // Package to work with time objects

	"github.com/aws/aws-lambda-go/events"                       // AWS Lambda package for Go events
	"github.com/aws/aws-lambda-go/lambda"                       // AWS Lambda package for Go
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin" // AWS Lambda custom runtime with Gin (web framework) support
	"github.com/gin-gonic/gin"                                  // Gin is a web framework for Go
)

// Declare global variable for ginLambda (AWS Lambda with Gin)
var ginLambda *ginadapter.GinLambda

// Define RequestBody structure to parse request JSON
type RequestBody struct {
	RequestId   string `json:"requestID"`
	RequestTime string `json:"requestTime"`
	Data        struct {
		Value1 int `json:"value1"`
		Value2 int `json:"value2"`
	} `json:"data"`
}

// Initialize Gin application and setup routes
func init() {
	log.Print("Gin cold start")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Route handler for POST requests to /api
	r.POST("/api", func(c *gin.Context) {
		// Bind request JSON to reqBody struct
		var reqBody RequestBody
		c.BindJSON(&reqBody)
		// Log request body
		log.Printf("%+v", reqBody)

		// Calculate sum of Value1 and Value2
		sum := reqBody.Data.Value1 + reqBody.Data.Value2

		// Send response JSON with requestId, current server time, and calculated sum
		c.JSON(200, gin.H{
			"requestId":   reqBody.RequestId,
			"requestTime": time.Now().Format(time.RFC3339),
			"data": gin.H{
				"sum": sum,
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
