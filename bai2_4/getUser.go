package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type Request struct {
	RequestID   string    `json:"requestId" validate:"required"`
	RequestTime time.Time `json:"requestTime"`
	Data        struct {
		Username string `json:"username" validate:"required"`
	} `json:"data" validate:"required"`
}

type Response struct {
	ResponseID      string    `json:"responseId"`
	ResponseTime    time.Time `json:"responseTime"`
	ResponseCode    string    `json:"responseCode"`
	ResponseMessage string    `json:"responseMessage"`
	Data            User      `json:"data"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse yêu cầu thành struct Request
	var req Request
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		response := Response{
			ResponseID:      req.RequestID,
			ResponseTime:    time.Now(),
			ResponseCode:    "400",
			ResponseMessage: "Lỗi khi phân tích yêu cầu",
		}
		responseBody, _ := json.Marshal(response)
		log.Println("Lỗi khi phân tích yêu cầu:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
		}, nil
	}

	// Kiểm tra tính hợp lệ của yêu cầu
	err = validateRequest(req)
	if err != nil {
		response := Response{
			ResponseID:      req.RequestID,
			ResponseTime:    time.Now(),
			ResponseCode:    "400",
			ResponseMessage: "Yêu cầu không hợp lệ",
		}
		responseBody, _ := json.Marshal(response)
		log.Println("Yêu cầu không hợp lệ:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
		}, nil
	}

	// Kết nối tới cơ sở dữ liệu PostgreSQL
	db, err := sql.Open("postgres", "postgres://thuypq123:t1xHgWF2DCAN@ep-empty-paper-000226.ap-southeast-1.aws.neon.tech/neondb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Lấy chi tiết người dùng
	var user User
	query := "SELECT username, name, phone FROM users WHERE username = $1"
	err = db.QueryRow(query, req.Data.Username).Scan(&user.Username, &user.Name, &user.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			response := Response{
				ResponseID:      req.RequestID,
				ResponseTime:    time.Now(),
				ResponseCode:    "400",
				ResponseMessage: "Không tìm thấy người dùng với username đã cho",
			}
			responseBody, _ := json.Marshal(response)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       string(responseBody),
			}, nil
		} else {
			log.Fatal(err)
		}
	}

	response := Response{
		ResponseID:      req.RequestID,
		ResponseTime:    time.Now(),
		ResponseCode:    "200",
		ResponseMessage: "Lấy chi tiết người dùng thành công",
		Data:            user,
	}

	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func validateRequest(req Request) error {
	validate := validator.New()
	err := validate.Struct(req)
	return err
}

func main() {
	lambda.Start(HandleRequest)
}
