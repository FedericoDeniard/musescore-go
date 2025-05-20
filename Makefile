APP_NAME=Musescore-Scrapper
BIN_DIR=dist

build:
	go build -o $(BIN_DIR)/$(APP_NAME) src/main.go

run:
	go run src/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME).exe ./src/main.go

clean: 
	rm -rf $(BIN_DIR)/*