
run:
	go build -o bin/web  cmd/web/web.go && ./bin/web

lint:
	golangci-lint run --fix