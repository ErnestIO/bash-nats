install:
	go install -v

build:
	go build -v ./...

lint:
	gometalinter

test:
	go test -v ./... --cover

deps:
	go get github.com/nats-io/nats
	go get github.com/ernestio/ernest-config-client
	go get golang.org/x/crypto/pbkdf2

dev-deps: deps
	go get github.com/alecthomas/gometalinter

clean:
	go clean

