package main

import (
	"bai3/common"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator"
)

// handler function to process the request and return a response
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		// Define variables to store request and response data
		req common.UserRequest
		res common.UserResponse
	)

	// Unmarshal the JSON body of the request into the UserRequest struct
	json.Unmarshal([]byte(request.Body), &req)

	// Log the unmarshalled request data for debugging purposes
	log.Println("Request:", req)

	// Validate the UserRequest struct based on predefined rules in the struct's tag
	validate := validator.New()
	err := validate.Struct(req)

	// If validation errors exist, log them and return an error response
	if err != nil {
		log.Println("- Error while parsing the request:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when parsing the request"}, nil
	}

	// Verify that the signature of the request is valid
	signatureIsValid := common.CheckSignature(req.Request, req.Data, req.Request.Signature, "golang")

	// If the signature is not valid, log it and return an error response
	if !signatureIsValid {
		log.Println("Invalid signature")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid signature"}, nil
	}

	// Convert phone number from string to integer
	numberPhone, err := strconv.Atoi(req.Data.Phone)

	// If there's an error in conversion, log it and return an error response
	if err != nil {
		fmt.Println("Error when converting phone to integer:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when handling phone number"}, nil
	}

	// Verify the validity of the phone number
	phoneIsValid, err := common.VerifyPhone(numberPhone)

	// If there's an error during verification or if the phone number is invalid, log it and return an error response
	if err != nil || !phoneIsValid {
		log.Println("Error when verifying phone:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when verifying phone"}, nil
	}

	// Create a user using the provided API
	UserResponse, err := common.CreateUser(req)

	// If there's an error during user creation or if the ResponseCode is not equal to "200", log it and return an error response
	if err != nil || UserResponse.ResponseCode != "200" {
		log.Println("Error when creating user:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when creating user"}, nil
	}
	// Marshal the response struct back into JSON
	responseBody, err := json.Marshal(res)

	if err != nil {
		log.Println("Error when creating response:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when creating response"}, nil
	}

	// Fill the response struct with necessary details
	res.ResponseCode = "00"
	res.Response.ResponseID = req.Request.RequestId
	res.ResponseMessage = "Success"
	res.Response.ResponseTime = req.Request.RequestTime

	// Return the successful response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

// The main function where the handler is registered and the lambda function starts
func main() {
	lambda.Start(handler)
}
