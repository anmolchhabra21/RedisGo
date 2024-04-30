run: build
	@./bin/redisgo
build: 
	@go build -o bin/redisgo .
