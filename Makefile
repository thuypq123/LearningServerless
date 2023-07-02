build:
	cd api2 && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/api2 Api2.go
	cd bai2_1 && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/bai2_1 createUSer.go
	cd bai3 && env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/bai3 bai3-api.go
	
deploy: build
	serverless deploy --aws-profile thuy --verbose