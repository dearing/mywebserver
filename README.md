# mywebserver

demo portable webserver with server side events, websockets and embedded assets

## about

This project serves to test and demostrate three other projects that are stand alone tools within Go's 1.24 module definition.

- https://github.com/dearing/go-cross-compile
- https://github.com/dearing/go-github-release
- https://github.com/dearing/go-live-reload


Where `go-live-reload` continuously builds and runs a set of projects, `go-cross-compile` builds binaries for combinations of architectures and operating systems and `go-github-release` publishes these artifacts on github. Developing a backend is what `go-live-reload` is mostly useful for so this project is a webserver that also serves to demonstrate and mess around with embedded filesystems, server side events, websockets, context, channels and templates.

## try out

Clone this project on a host with Go 1.24+ and the tools will get pulled and compiled on demand.

```
git clone https://github.com/dearing/mywebserver.git
cd mywebserver
go tool go-live-reload --config-file .config/go-live-reload.json
```

## go tool demostration
```
go tool go-live-reload --config-file .config/go-live-reload.json --build-groups www
go tool go-cross-compile --config-file .config/go-cross-compile.json
go tool go-github-release --github-owner dearing --github-repo mywebserver --tag-name v1.0.1
```
## webserver usage

```
Usage: mywebserver [options]

Demo portable webserver with server side events and embedded assets.

- https://github.com/dearing/mywebserver

Options:

  -bind string
        bind address (default ":8080")
  -sse-duration duration
        SSE ticker duration (default 1s)
  -version
        print version information
```
