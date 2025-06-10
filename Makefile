build:
	go build -o go-q cmd/server/server.go
	codesign --sign - --force --deep ./go-q
run:
	./go-q
start:
	make build
	make run
