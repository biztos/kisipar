default:
	go build -o kisipar cmd/kisipar/main.go

test:
	go test ./...
