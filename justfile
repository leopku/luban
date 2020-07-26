dev:
  watchexec -e go -r -- go run main.go

wire:
  wire ./...

dep:
  go get -u github.com/go-swagger/go-swagger/cmd/swagger

test:
  # go run main.go mysql
  go test $(go list ./... | grep -v generated)