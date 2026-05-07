module example

go 1.25.9

require mycache v0.0.0

require (
	github.com/golang/protobuf v1.5.4 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace mycache => ./my-cache
