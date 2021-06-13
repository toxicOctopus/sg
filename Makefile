run-web:
	go build -o bin/web  cmd/web/web.go && ./bin/web

win-run-web:
	go build -o bin\web.exe  cmd\web\web.go & bin\web.exe

lint:
	golangci-lint run

generate-config:
	go build -o bin/cli  cmd/cli/cli.go && ./bin/cli generate-config

win-generate-config:
	go build -o bin\cli.exe  cmd\cli\cli.go && bin\cli.exe generate-config

lin-centrifugo:
	sudo docker run --ulimit nofile=65536:65536 -v /home/toxocto/Projects/sg/config/centrifugo.json:/centrifugo/centrifugo.json -p 8000:8000 centrifugo/centrifugo centrifugo -c centrifugo.json

win-centrifugo:
	docker run --ulimit nofile=65536:65536 -v C:\Users\to-pc\Projects\sg\config\centrifugo.json:/centrifugo/centrifugo.json -p 8000:8000 centrifugo/centrifugo centrifugo -c centrifugo.json

test:
	go test ./... -v