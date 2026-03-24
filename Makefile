.PHONY: build test lint clean

build:
	go build -o dbx-cli .

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -f dbx-cli
