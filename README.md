# SimpleHTTPServer-golang

## Usage

### Help
```
$ ./webserver -help
web subcommand
  -bind string
    	Bind Port (default "8080")
  -certificate string
    	HTTPS certificate : openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365
  -gzip
    	Enables gzip/zlib compression (default true)
  -help
    	Print usage
  -key string
    	HTTPS Key : openssl genrsa -out server.key 2048
  -password string
    	Password for basic auth, default: notsecure (default "notsecure")
  -private string
    	Private folder with basic auth, default /tmp/SimpleHTTPServer-golang/src/bin/private (default "private")
  -public string
    	Default /home/florian/SimpleHTTPServer-golang/src/bin public folder
  -root string
    	Root folder (default "/tmp/SimpleHTTPServer-golang/src/bin")
  -tls
    	Enables HTTPS
  -username string
    	Username for basic auth, default: admin (default "admin")

run subcommand
Usage:
./webserver-linux-amd64 run <binary> <args>

Packaged Binaries:
```

### Web server over HTTP
```
$ ./webserver
Sharing /tmp/ on 8080 ...
Sharing /tmp/private on 8080 ...
```

### Web server over HTTPS
```
$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
..........................................+++++
.................................................................................................................+++++
e is 65537 (0x010001)

$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:
Email Address []:

$ ./webserver web -tls -key server.key -cert server.crt
Sharing /tmp/ on 8080 ...
Sharing /tmp/private on 8080 ...
```

### Web server using private directory
```
$ ./webserver-linux-amd64 web -private ThePrivateFolder -username nodauf -password nodauf
Sharing /tmp/ on 8080 ...
Sharing /tmp/ThePrivateFolder on 8080 ...
```

### Run embedded binary
```
C:/Users/Nodauf>.\webserver.exe run Rubeus.exe kerberoast 
```

## Features

* HTTPS
* Directory listing
* Define a private folder with basic authentication
* Upload file
* Download file as an encrypted zip (password: infected)
* Download folder with a zip
* Embedded files
* Run embedded binary written in C# (only available on Windows)

## Todo

* Ability to execute embedded binary in C# and C
* Implement (go-mimikatz)[https://github.com/vyrus001/go-mimikatz]
* Add feature for search and replace in embedded files (for fill the IP address for example)
* JS/CSS menu to give command line in powershell, some gtfobins, curl, wget to download and execute 
* Create a folder for the browser
