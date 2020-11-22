build:
	go build -o bin/webserver webserver.go

run:
	go run webserver.go

compile:
	# 32-Bit Systems
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
    
	# 64-Bit
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 main.go

