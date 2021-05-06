run-web:
	go build -o bin/web  cmd/web/web.go && ./bin/web

lint:
	golangci-lint run

generate-config:
	go build -o bin/cli  cmd/cli/cli.go && ./bin/cli generate-config

compose-up:
	sudo docker run --ulimit nofile=65536:65536 -v /home/toxocto/Projects/sg/config/centrifugo.json:/centrifugo/centrifugo.json -p 8000:8000 centrifugo/centrifugo centrifugo -c centrifugo.json

test:
	go test ./... -v