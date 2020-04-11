run:
	go build -o bin/web  cmd/web/web.go && ./bin/web

lint:
	golangci-lint run

generate-config:
	go build -o bin/cli  cmd/cli/cli.go && ./bin/cli generate-config