package main

// import lambda and write handler function
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

// write function handler for lambda
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		req common.UserRequest
		res common.UserResponse
	)
	// get request body
	json.Unmarshal([]byte(request.Body), &req)
	// log request
	log.Println("Request:", req)
	// validate request with validator struct
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		log.Println("-Lỗi khi phân tích yêu cầu:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Lỗi khi phân tích yêu cầu"}, nil
	}
	// verify signature with secret key with hmac CheckSignature function
	signatureIsValid := common.CheckSignature(req.Request, req.Data, req.Request.Signature, "golang")
	if !signatureIsValid {
		log.Println("Invalid signature")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid signature"}, nil
	}
	// verify phone with VerifyPhone function
	numberPhone, err := strconv.Atoi(req.Data.Phone)
	if err != nil {
		fmt.Println("Error when convert phone to int:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when handle number phone"}, nil
	}
	phoneIsValid, err := common.VerifyPhone(numberPhone)
	if err != nil {
		log.Println("Error when verify phone:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when verify phone"}, nil
	}
	if !phoneIsValid {
		log.Println("Phone is invalid")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Phone is invalid"}, nil
	}
	// create user with Api
	UserResponse, err := common.CreateUser(req)
	if err != nil {
		log.Println("Error when create user:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error when create user"}, nil
	}
	if UserResponse.ResponseCode != "200" {
		log.Println("Create user failed")
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Create user failed"}, nil
	}
	// create response
	res.ResponseCode = "00"
	res.Response.ResponseID = req.Request.RequestId
	res.ResponseMessage = "Success"
	res.Response.ResponseTime = req.Request.RequestTime

	// create response data
	responseBody, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}
