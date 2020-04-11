run-web:
	go build -o bin/web  cmd/web/web.go && ./bin/web

lint:
	golangci-lint run

generate-config:
	go build -o bin/cli  cmd/cli/cli.go && ./bin/cli generate-config

compose-up:
	sudo docker run --ulimit nofile=65536:65536 -v /home/toxocto/Downloads/centrifugo_2.3.1_linux_amd64/:/centrifugo -p 8000:8000 centrifugo/centrifugo centrifugo -c config/centrifugo.json