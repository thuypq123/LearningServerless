build:
	 go build -ldflags="-s -w" -o bin/bai2_1 bai2_1/CreateUSer.go
	 go build -ldflags="-s -w" -o bin/bai2_2 bai2_2/UpdateUser.go
	
deploy: build
	serverless deploy --verbose


	# $Env:GOOS = "linux"; $Env:GOARCH = "amd64";