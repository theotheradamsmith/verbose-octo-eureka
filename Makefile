BIN = verify
VERIFY_SRC = src/*.go

all: $(BIN)

docker: # Build a binary in docker and copy it to the local filesystem
	docker build -t go-build .
	docker container create --name temp go-build
	docker container cp temp:/go/src/app/verify ./
	docker container rm temp

verify: $(VERIFY_SRC)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@ $^ 

clean:
	rm -f verify
	go clean