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
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type Request struct {
	RequestID   string    `json:"requestId" validate:"required"`
	RequestTime time.Time `json:"requestTime"`
	Data        User      `json:"data" validate:"required"`
}

type Response struct {
	ResponseID      string    `json:"responseId"`
	ResponseTime    time.Time `json:"responseTime"`
	ResponseCode    string    `json:"responseCode"`
	ResponseMessage string    `json:"responseMessage"`
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

	// Kiểm tra xem người dùng có tồn tại trong cơ sở dữ liệu hay không
	var existingUsername string
	err = db.QueryRow("SELECT username FROM users WHERE username = $1", req.Data.Username).Scan(&existingUsername)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if existingUsername == "" {
		response := Response{
			ResponseID:      req.RequestID,
			ResponseTime:    time.Now(),
			ResponseCode:    "400",
			ResponseMessage: "Người dùng không tồn tại",
		}
		responseBody, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
		}, nil
	}

	// Cập nhật thông tin người dùng
	query := "UPDATE users SET name = $1, phone = $2 WHERE username = $3"
	_, err = db.Exec(query, req.Data.Name, req.Data.Phone, req.Data.Username)
	if err != nil {
		log.Fatal(err)
	}

	response := Response{
		ResponseID:      req.RequestID,
		ResponseTime:    time.Now(),
		ResponseCode:    "200",
		ResponseMessage: "Cập nhật người dùng thành công",
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
