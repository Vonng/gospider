upload:
	GOOS=linux GOARCH=amd64 go build -o wdj_app wdj_app.go
	oss put wdj_app

download:
	oss get wdj_app

linux:
	GOOS=linux GOARCH=amd64 go build -o wdj_app wdj_app.go

mac:
	GOOS=darwin GOARCH=amd64 go build -o wdj_app wdj_app.go

build:
	go build -o wdj_app wdj_app.go

clean:
	rm -rf wdj_app

