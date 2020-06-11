dev:
  watchexec -e go -r -- go run main.go

wire
  wire ./...

dep
  go get -u github.com/go-swagger/go-swagger/cmd/swagger
