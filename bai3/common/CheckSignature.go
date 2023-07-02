package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func CheckSignature(request Request, data User, signature, secretKey string) bool {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secretKey))

	// hash = requestId+phone+username+secretKey
	requestId := request.RequestId
	phone := data.Phone
	username := data.Username

	// print(requestId + phone + username + secretKey)
	log.Println(requestId + phone + username + secretKey)
	h.Write([]byte(requestId + phone + username + secretKey))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))
	log.Println("sha:", sha)

	return sha == signature
}
