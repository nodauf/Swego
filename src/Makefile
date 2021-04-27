build: assetsDir
	go build -ldflags="-s -w" -o bin/webserver webserver.go
	rm controllers/rice-box.go

run:
	go run webserver.go

assetsDir:
	go get github.com/GeertJohan/go.rice
	go get github.com/GeertJohan/go.rice/rice
	$(GOPATH)/bin/rice embed-go -i ./controllers

compileLinux: assetsDir
	# 32-Bit Systems
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o bin/webserver-linux-386 webserver.go
    
	# 64-Bit
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/webserver-linux-amd64 webserver.go
	# Remove rice file
	rm controllers/rice-box.go

compileWindows: assetsDir
	# 32-Bit Systems
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o bin/webserver-windows-386.exe webserver.go
    
	# 64-Bit
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/webserver-windows-amd64.exe webserver.go
	# Remove rice file
	rm controllers/rice-box.go
