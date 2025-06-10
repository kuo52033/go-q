build:
	go build -o ./tmp/go-q cmd/server/server.go
	codesign --sign - --force --deep ./tmp/go-q
run:
	./tmp/go-q
start:
	make build
	make run
