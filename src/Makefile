build:
	go build -o bin/webserver webserver.go

run:
	go run webserver.go

assetsDir:
	go get github.com/GeertJohan/go.rice
	go get github.com/GeertJohan/go.rice/rice
	$(GOPATH)/bin/rice embed-go -i ./controllers

compile: assetsDir
	# 32-Bit Systems
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o bin/webserver-linux-386 webserver.go
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/webserver-windows-386.exe webserver.go
    
	# 64-Bit
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/webserver-linux-amd64 webserver.go
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/webserver-windows-amd64.exe webserver.go
	# Remove rice file
	rm controllers/rice-box.go
