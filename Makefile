default: build

build:
	go build -o deck_pic_generator main.go
	GOOS=windows GOARCH=amd64 go build -o deck_pic_generator.exe main.go
