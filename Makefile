build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main.go 
serverless:
	sls deploy --aws-profile default

	$Env:GOOS = "linux"; $Env:GOARCH = "amd64"