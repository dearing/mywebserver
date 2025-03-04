module github.com/dearing/mywebserver

go 1.24.0

tool (
	github.com/dearing/go-cross-compile
	github.com/dearing/go-github-release
	github.com/dearing/go-live-reload
	github.com/dearing/mywebserver
)

require golang.org/x/net v0.35.0

require (
	github.com/dearing/go-cross-compile v1.0.8 // indirect
	github.com/dearing/go-github-release v1.0.3 // indirect
	github.com/dearing/go-live-reload v0.4.0 // indirect
)
