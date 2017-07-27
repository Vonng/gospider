upload:
	GOOS=linux GOARCH=amd64 go build -o ios_app ios_app.go
	oss put ios_app

download:
	oss get ios_app

linux:
	GOOS=linux GOARCH=amd64 go build -o ios_app ios_app.go

mac:
	GOOS=darwin GOARCH=amd64 go build -o ios_app ios_app.go

build:
	go build -o ios_app ios_app.go

clean:
	rm -rf ios_app

