BIN = verify
VERIFY_SRC = main.go

all: $(BIN)

verify: $(VERIFY_SRC)
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@ $^ 

clean:
	rm -f verify
	go clean