package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
	// err = validateRequest(req)
	// if err != nil {
	// 	response := Response{
	// 		ResponseID:      req.RequestID,
	// 		ResponseTime:    time.Now(),
	// 		ResponseCode:    "400",
	// 		ResponseMessage: "Yêu cầu không hợp lệ",
	// 	}
	// 	responseBody, _ := json.Marshal(response)
	// 	log.Println("Yêu cầu không hợp lệ:", err)
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 400,
	// 		Body:       string(responseBody),
	// 	}, nil
	// }

	// Kết nối tới cơ sở dữ liệu PostgreSQL
	db, err := sql.Open("postgres", "postgres://thuypq123:t1xHgWF2DCAN@ep-empty-paper-000226.ap-southeast-1.aws.neon.tech/neondb")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// kiểm tra trùng username
	var username string
	queryCheck := `SELECT username FROM users WHERE username = $1`
	err = db.QueryRow(queryCheck, req.Data.Username).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Không có user nào có username này")
		} else {
			log.Fatal(err)
		}
	}
	if username != "" {
		response := Response{
			ResponseID:      req.RequestID,
			ResponseTime:    time.Now(),
			ResponseCode:    "400",
			ResponseMessage: "Username đã tồn tại",
		}
		responseBody, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
		}, nil
	}
	// Tạo một user mới từ dữ liệu yêu cầu
	user := User{
		Username: req.Data.Username,
		Name:     req.Data.Name,
		Phone:    req.Data.Phone,
	}

	// Thực hiện truy vấn INSERT để thêm user vào cơ sở dữ liệu
	query := `INSERT INTO users (username, name, phone) VALUES ($1, $2, $3) RETURNING id`
	err = db.QueryRow(query, user.Username, user.Name, user.Phone).Scan(&user.ID)
	if err != nil {
		fmt.Println("Lỗi khi thêm user:", err)
		log.Fatal(err)
	}

	// In thông tin user đã được thêm vào cơ sở dữ liệu
	fmt.Printf("Thêm user thành công. ID: %d, Username: %s, Name: %s, Phone: %s\n", user.ID, user.Username, user.Name, user.Phone)

	response := Response{
		ResponseID:      req.RequestID,
		ResponseTime:    time.Now(),
		ResponseCode:    "200",
		ResponseMessage: "Thêm user thành công",
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
	// mock call lambda
	// request := Request{
	// 	RequestID:   "123",
	// 	RequestTime: time.Now(),
	// 	Data: User{
	// Username: "thuypq123",
	// Name:     "Phan Quang Thuy",
	// Phone:    "0123456789",
	// 	},
	// }
	// requestBody, _ := json.Marshal(request)
	// response, _ := HandleRequest(context.Background(), events.APIGatewayProxyRequest{
	// 	Body: string(requestBody),
	// })
	// fmt.Println(response)
	lambda.Start(HandleRequest)
}
